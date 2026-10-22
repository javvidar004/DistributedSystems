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
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

func CreateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.CreateUserRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user payload"})
			return
		}

		if _, err := getUserByEmail(db, payload.Email); err == nil {
			c.JSON(http.StatusConflict, gin.H{"message": "email already exists"})
			return
		}

		user := models.User{
			Email:        payload.Email,
			Name:         payload.Name,
			LastName:     payload.LastName,
			WorkPosition: payload.WorkPosition,
			Salary:       payload.Salary,
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user", "error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "user created successfully", "user": user})
	}
}

func UpdateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.UpdateUserRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_user", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user payload"})
			return
		}

		user, err := getUserByEmail(db, payload.Email)
		if err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_user", "failed")
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		user.Name = payload.Name
		user.LastName = payload.LastName
		user.WorkPosition = payload.WorkPosition
		user.Salary = payload.Salary
		if err := db.Save(&user).Error; err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_user", "failed")
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not update user", "error": err.Error()})
			return
		}

		newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_user", "success")
		c.JSON(http.StatusOK, gin.H{"message": "user updated successfully", "user": user})
	}
}

func DeleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == "" {
			username = c.Query("username")
		}
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is required"})
			return
		}
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "username is required"})
			return
		}
		result := db.Where("username = ?", username).Delete(&models.User{})
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete user", "error": result.Error.Error()})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
	}
}

func GetUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "userID is required"})
			return
		}

		user, err := getUserByID(db, userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user found successfully", "user": user})
	}
}

func GetUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		result := db.Select("id", "email", "name", "last_name", "work_position", "salary").Find(&users)
		//user := models.User{}
		//result := db.Select("email", "name", "last_name", "work_position", "salary").Find(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not get users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users}) //gin.H{"message": "users found successfully", "users": users})
	}
}

func getUserByID(db *gorm.DB, userID string) (models.User, error) {
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
}

func getUserByEmail(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, err
		}
		return models.User{}, err
	}
	return user, nil
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

func HeartbeatHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		instanceID := os.Getenv("INSTANCE_ID")
		c.JSON(http.StatusOK, gin.H{"message": "heartbeat received", "instanceID": instanceID})
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
		ServiceName: "users",
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

}

func RegisterWithConsul(consulAddr string, instanceID string) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "users-service",
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
