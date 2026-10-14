package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

func RegisterWithConsul(consulAddr string, instanceID string) error {
	payload := map[string]any{
		"ID":      instanceID,
		"Name":    "logs-service",
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
