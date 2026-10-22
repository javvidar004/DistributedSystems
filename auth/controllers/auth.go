package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.LoginRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "login", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid login payload"})
			return
		}

		user, err := getUserByUsername(db, payload.Email)
		if err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "login", "failed")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
			return
		}

		//if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		//	c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		//	return
		//}

		if user.PasswordHash != payload.Password {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "login", "failed")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
			return
		}

		newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "login", "success")
		c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user})
	}
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.RegisterRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid register payload"})
			return
		}

		if payload.Email == "" || payload.Password == "" {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "email and password are required"})
			return
		}

		if _, err := getUserByUsername(db, payload.Email); err == nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusConflict, gin.H{"message": "username already exists"})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not hash password", "error": err.Error()})
			return
		}

		uuid, err := uuid.NewRandom()
		if err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate UUID", "error": err.Error()})
			return
		}

		user := models.User{
			ID:           uuid.String(),
			Email:        payload.Email,
			PasswordHash: strings.TrimSpace(string(hash)),
			IsActive:     true,
			Role:         "user",
			Group:        "default",
		}

		if err := db.Create(&user).Error; err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user", "error": err.Error()})
			return
		}

		newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "success")
		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
	}
}

func UpdatePasswordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.UpdatePasswordRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password payload"})
			return
		}

		if payload.Email == "" || payload.NewPassword == "" {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "email and new password are required"})
			return
		}

		user, err := getUserByUsername(db, payload.Email)
		if err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "failed")
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		if user.PasswordHash != payload.OldPassword {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "failed")
			c.JSON(http.StatusUnauthorized, gin.H{"message": "old password is incorrect"})
			return
		}

		user.PasswordHash = payload.NewPassword
		if err := db.Save(&user).Error; err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update password", "error": err.Error()})
			return
		}

		newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "success")
		c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
	}
}

func DeleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			userID = c.Query("userID")
		}
		if userID == "" {
			newLogRequest(time.Now().Format(time.RFC3339), "", "delete_user", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "userID is required"})
			return
		}

		result := db.Where("id = ?", userID).Delete(&models.User{})
		if result.Error != nil {
			newLogRequest(time.Now().Format(time.RFC3339), userID, "delete_user", "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete user", "error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			newLogRequest(time.Now().Format(time.RFC3339), userID, "delete_user", "failed")
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		newLogRequest(time.Now().Format(time.RFC3339), userID, "delete_user", "success")
		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	}
}

func getUserByUsername(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
}

//func withInstanceHeader(next http.HandlerFunc) http.HandlerFunc {
//
//	var instanceID string
//	if v := os.Getenv("INSTANCE_ID"); v != "" {
//		instanceID = v
//	}
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("X-Instance-Id", instanceID)
//		next(w, r)
//	}
//}

func HeartbeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		instanceID := os.Getenv("INSTANCE_ID")
		c.JSON(http.StatusOK, gin.H{"message": "heartbeat received", "instanceID": instanceID})
	}
}

func WithInstanceHeader(next gin.HandlerFunc) gin.HandlerFunc {

	var instanceID string
	if v := os.Getenv("INSTANCE_ID"); v != "" {
		instanceID = v
	}
	return func(c *gin.Context) {
		c.Header("X-Instance-Id", instanceID)
		next(c)
	}
}

func newLogRequest(timestamp, username, action, status string) {
	if err := newLogRequestGRPC(timestamp, username, action, status); err == nil {
		return
	}

	newLogRequestHTTP(timestamp, username, action, status)
}

func newLogRequestGRPC(timestamp, username, action, status string) error {
	addr := os.Getenv("LOGS_GRPC_ADDR")
	if addr == "" {
		addr = "logsbe1:9090"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewLogServiceClient(conn)
	_, err = client.CreateLog(ctx, &pb.CreateLogRequest{
		ServiceName: "auth",
		Action:      action,
		UserId:      username,
		Details:     "status=" + status + ";timestamp=" + timestamp,
	})
	if err != nil {
		return err
	}

	return nil
}

func newLogRequestHTTP(timestamp, username, action, status string) {

	body := conectors.NewLogRequest{
		Timestamp: timestamp,
		Username:  username,
		Action:    action,
		Status:    status,
	}

	// Send request to LogLB service
	url := "http://" + os.Getenv("MIDDLEWARE_HOST") + ":8080/logs/log"

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling log request: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating log request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending log request: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Log request sent: %s, response status: %s", action, resp.Status)

}

func RegisterWithConsul(consulAddr string, instanceID string) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "auth-service",
		"Address": instanceID, // resolvable via Docker's embedded DNS
		"Port":    8080,
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:8080/heartbeat", instanceID),
			"Interval":                       "5s",
			"Timeout":                        "2s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= 10; attempt++ {
		req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/register", bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Printf("registered %s with Consul at %s", instanceID, consulAddr)
				return nil
			}
			lastErr = fmt.Errorf("consul returned status %d", resp.StatusCode)
		} else {
			lastErr = err
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return fmt.Errorf("giving up registering with consul: %w", lastErr)
}

func DeregisterFromConsul(consulAddr string, instanceID string) {
	req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/deregister/"+instanceID, nil)
	if err != nil {
		log.Printf("failed to build deregister request: %v", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("failed to deregister %s from consul: %v", instanceID, err)
		return
	}
	resp.Body.Close()
	log.Printf("deregistered %s from consul", instanceID)
}
