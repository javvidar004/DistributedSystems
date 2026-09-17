package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
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
		result := db.Select("email", "name", "last_name", "work_position", "salary").Find(&users)
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

	body := conectors.NewLogRequest{
		Timestamp: timestamp,
		Username:  username,
		Action:    action,
		Status:    status,
	}

	// Send request to LogLB service
	url := "http://" + os.Getenv("LOG_LB_HOST") + ":8080/log"

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
