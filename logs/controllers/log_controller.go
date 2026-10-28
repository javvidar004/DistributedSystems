package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"proyecto.com/m/conectors"
	"proyecto.com/m/service"
)

func GetLogHandler(logsSvc *service.LogsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := logsSvc.GetLogs(c.Request.Context())
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"logs": res.Logs})
	}
}

func NewLogHandler(logsSvc *service.LogsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.NewLogRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid log payload"})
			return
		}

		res, err := logsSvc.NewLog(c.Request.Context(), service.NewLogInput{
			Timestamp: payload.Timestamp,
			Username:  payload.Username,
			Action:    payload.Action,
			Status:    payload.Status,
		})
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": res.Message, "log": res.Log})
	}
}

func mapHTTPError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	switch {
	case errors.Is(err, service.ErrInvalidPayload):
		return http.StatusBadRequest, "invalid log payload"
	default:
		return http.StatusInternalServerError, "internal server error"
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
	httpPort := 8080
	if v := os.Getenv("HTTP_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			httpPort = p
		}
	}
	grpcPort := 50051
	if v := os.Getenv("GRPC_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			grpcPort = p
		}
	}

	httpPayload := map[string]any{
		"ID":      instanceID,
		"Name":    "logs-service",
		"Address": instanceID, // resolvable via Docker's embedded DNS
		"Port":    httpPort,
		"Check": map[string]any{
			"HTTP":                           fmt.Sprintf("http://%s:%d/heartbeat", instanceID, httpPort),
			"Interval":                       "5s",
			"Timeout":                        "2s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}

	grpcPayload := map[string]any{
		"ID":      instanceID + "-grpc",
		"Name":    "logs-grpc-service",
		"Address": instanceID,
		"Port":    grpcPort,
		"Check": map[string]any{
			"TCP":                            fmt.Sprintf("%s:%d", instanceID, grpcPort),
			"Interval":                       "5s",
			"Timeout":                        "2s",
			"DeregisterCriticalServiceAfter": "1m",
		},
	}

	httpBody, err := json.Marshal(httpPayload)
	if err != nil {
		return err
	}
	grpcBody, err := json.Marshal(grpcPayload)
	if err != nil {
		return err
	}

	var lastErr error
	bodies := []struct {
		body []byte
		name string
	}{
		{body: httpBody, name: "http"},
		{body: grpcBody, name: "grpc"},
	}

	for _, entry := range bodies {
		registered := false
		for attempt := 1; attempt <= 10; attempt++ {
			req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/register", bytes.NewReader(entry.body))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					log.Printf("registered %s %s with Consul at %s", instanceID, entry.name, consulAddr)
					registered = true
					break
				}
				lastErr = fmt.Errorf("consul returned status %d", resp.StatusCode)
			} else {
				lastErr = err
			}
			time.Sleep(time.Duration(attempt) * time.Second)
		}
		if !registered {
			return fmt.Errorf("giving up registering %s endpoint with consul: %w", entry.name, lastErr)
		}
	}

	return nil
}

func DeregisterFromConsul(consulAddr string, instanceID string) {
	serviceIDs := []string{instanceID, instanceID + "-grpc"}
	for _, svcID := range serviceIDs {
		req, err := http.NewRequest(http.MethodPut, consulAddr+"/v1/agent/service/deregister/"+svcID, nil)
		if err != nil {
			log.Printf("failed to build deregister request for %s: %v", svcID, err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("failed to deregister %s from consul: %v", svcID, err)
			continue
		}
		resp.Body.Close()
		log.Printf("deregistered %s from consul", svcID)
	}
}
