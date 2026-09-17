package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"proyecto.com/m/conectors"
	"proyecto.com/m/models"
)

func GetLogHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logs []models.Log
		if err := db.Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch logs", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

func NewLogHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.NewLogRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid log payload"})
			return
		}

		log := models.Log{
			ID:        uuid.New().String(),
			Timestamp: payload.Timestamp,
			Username:  payload.Username,
			Action:    payload.Action,
			Status:    payload.Status,
		}

		if err := db.Create(&log).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create log", "error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "log created successfully", "log": log})
	}
}

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
