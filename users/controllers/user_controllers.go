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

func CreateUserHandler(usersSvc *service.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.CreateUserRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user payload"})
			return
		}

		res, err := usersSvc.CreateUser(c.Request.Context(), service.CreateUserInput{
			Email:        payload.Email,
			Name:         payload.Name,
			LastName:     payload.LastName,
			WorkPosition: payload.WorkPosition,
			Salary:       payload.Salary,
		})
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": res.Message, "user": res.User})
	}
}

func UpdateUserHandler(usersSvc *service.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.UpdateUserRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_user", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user payload"})
			return
		}

		res, err := usersSvc.UpdateUser(c.Request.Context(), service.UpdateUserInput{
			Email:        payload.Email,
			Name:         payload.Name,
			LastName:     payload.LastName,
			WorkPosition: payload.WorkPosition,
			Salary:       payload.Salary,
		})
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": res.Message, "user": res.User})
	}
}

func DeleteUserHandler(usersSvc *service.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			userID = c.Query("userID")
		}
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "userID is required"})
			return
		}

		res, err := usersSvc.DeleteUser(c.Request.Context(), userID)
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	}
}

func GetUserHandler(usersSvc *service.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "userID is required"})
			return
		}

		res, err := usersSvc.GetUser(c.Request.Context(), userID)
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": res.Message, "user": res.User})
	}
}

func GetUsersHandler(usersSvc *service.UsersService) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := usersSvc.GetUsers(c.Request.Context())
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": res.Users})
	}
}

func mapHTTPError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	switch {
	case errors.Is(err, service.ErrRequiredFields):
		return http.StatusBadRequest, "required fields are missing"
	case errors.Is(err, service.ErrUserAlreadyExists):
		return http.StatusConflict, "email already exists"
	case errors.Is(err, service.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	default:
		return http.StatusInternalServerError, "internal server error"
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
		"Name":    "users-service",
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
		"Name":    "users-grpc-service",
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
