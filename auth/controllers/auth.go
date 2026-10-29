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
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"

	"proyecto.com/m/conectors"
	"proyecto.com/m/service"
)

func LoginHandler(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.LoginRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "login", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid login payload"})
			return
		}

		res, err := authSvc.Login(c.Request.Context(), payload.Email, payload.Password)
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		tokens, err := authSvc.CreateJWT(res.User)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create auth token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       res.Message,
			"user":          res.User,
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
		})
	}
}

func ValidateJWTHandler(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if token == "" {
			var payload conectors.ValidateJWTRequest
			if err := c.ShouldBindJSON(&payload); err == nil {
				token = payload.Token
			}
		}

		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "token is required"})
			return
		}

		claims, err := authSvc.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "token is valid",
			"user_id": claims.UserID,
			"role":    claims.Role,
			"group":   claims.Group,
		})
	}
}

func RegisterHandler(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.RegisterRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "register", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid register payload"})
			return
		}

		res, err := authSvc.Register(c.Request.Context(), service.RegisterInput{
			Name:         payload.Name,
			LastName:     payload.LastName,
			Email:        payload.Email,
			WorkPosition: payload.WorkPosition,
			Salary:       payload.Salary,
			Password:     payload.Password,
		})
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": res.Message, "user": res.User})
	}
}

func UpdatePasswordHandler(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload conectors.UpdatePasswordRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			newLogRequest(time.Now().Format(time.RFC3339), payload.Email, "update_password", "failed")
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid password payload"})
			return
		}

		res, err := authSvc.UpdatePassword(c.Request.Context(), payload.Email, payload.OldPassword, payload.NewPassword)
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	}
}

func DeleteUserHandler(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userID")
		if userID == "" {
			userID = c.Query("userID")
		}
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "userID is required"})
			return
		}

		res, err := authSvc.DeleteUser(c.Request.Context(), userID)
		if err != nil {
			statusCode, msg := mapHTTPError(err)
			c.JSON(statusCode, gin.H{"message": msg})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": res.Message})
	}
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

func mapHTTPError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	switch {
	case errors.Is(err, service.ErrInvalidPayload):
		return http.StatusBadRequest, "invalid payload"
	case errors.Is(err, service.ErrRequiredFields):
		return http.StatusBadRequest, "email and password are required"
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized, "invalid username or password"
	case errors.Is(err, service.ErrUserAlreadyExists):
		return http.StatusConflict, "username already exists"
	case errors.Is(err, service.ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	case errors.Is(err, service.ErrOldPasswordMismatch):
		return http.StatusUnauthorized, "old password is incorrect"
	case errors.Is(err, service.ErrInvalidToken):
		return http.StatusUnauthorized, "invalid token"
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

func HeartbeatGRPCContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 3*time.Second)
}

func newLogRequest(timestamp, username, action, status string) {
	payload := map[string]string{
		"timestamp": timestamp,
		"username":  username,
		"action":    action,
		"status":    status,
	}

	message, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling kafka log payload: %v", err)
		return
	}

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_8_0_0
	config.Metadata.AllowAutoTopicCreation = true

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		log.Printf("error connecting to kafka: %v", err)
		return
	}
	defer producer.Close()

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "logs_topic",
		Value: sarama.ByteEncoder(message),
	})
	if err != nil {
		log.Printf("error sending log to kafka: %v", err)
		return
	}

	log.Printf("log request sent to kafka topic logs_topic: %s", action)
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
		"Name":    "auth-service",
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
		"Name":    "auth-grpc-service",
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
