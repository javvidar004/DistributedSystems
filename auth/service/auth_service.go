package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"proyecto.com/m/models"
)

var (
	ErrInvalidPayload      = errors.New("invalid payload")
	ErrRequiredFields      = errors.New("required fields are missing")
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("username already exists")
	ErrOldPasswordMismatch = errors.New("old password is incorrect")
)

type AuthService struct {
	db *gorm.DB
}

type RegisterInput struct {
	Name         string
	LastName     string
	Email        string
	WorkPosition string
	Salary       float64
	Password     string
}

type LoginResult struct {
	Message string
	User    models.User
}

type RegisterResult struct {
	Message string
	User    models.User
}

type MessageResult struct {
	Message string
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	if email == "" || password == "" {
		return LoginResult{}, ErrRequiredFields
	}

	user, err := s.getUserByEmail(ctx, email)
	if err != nil {
		newLogRequest(ctx, email, "login", "failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if user.PasswordHash != password {
		newLogRequest(ctx, email, "login", "failed")
		return LoginResult{}, ErrInvalidCredentials
	}

	newLogRequest(ctx, email, "login", "success")
	return LoginResult{Message: "login successful", User: user}, nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (RegisterResult, error) {
	if input.Email == "" || input.Password == "" {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, ErrRequiredFields
	}

	if _, err := s.getUserByEmail(ctx, input.Email); err == nil {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, fmt.Errorf("could not hash password: %w", err)
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, fmt.Errorf("could not generate UUID: %w", err)
	}

	user := models.User{
		ID:           uid.String(),
		Email:        input.Email,
		PasswordHash: strings.TrimSpace(string(hash)),
		IsActive:     true,
		Role:         "user",
		Group:        "default",
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, fmt.Errorf("could not create user: %w", err)
	}

	newLogRequest(ctx, input.Email, "register", "success")
	return RegisterResult{Message: "user created successfully", User: user}, nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, email, oldPassword, newPassword string) (MessageResult, error) {
	if email == "" || newPassword == "" {
		newLogRequest(ctx, email, "update_password", "failed")
		return MessageResult{}, ErrRequiredFields
	}

	user, err := s.getUserByEmail(ctx, email)
	if err != nil {
		newLogRequest(ctx, email, "update_password", "failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return MessageResult{}, ErrUserNotFound
		}
		return MessageResult{}, err
	}

	if user.PasswordHash != oldPassword {
		newLogRequest(ctx, email, "update_password", "failed")
		return MessageResult{}, ErrOldPasswordMismatch
	}

	user.PasswordHash = newPassword
	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		newLogRequest(ctx, email, "update_password", "failed")
		return MessageResult{}, fmt.Errorf("could not update password: %w", err)
	}

	newLogRequest(ctx, email, "update_password", "success")
	return MessageResult{Message: "password updated successfully"}, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userID string) (MessageResult, error) {
	if userID == "" {
		newLogRequest(ctx, "", "delete_user", "failed")
		return MessageResult{}, ErrRequiredFields
	}

	result := s.db.WithContext(ctx).Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		newLogRequest(ctx, userID, "delete_user", "failed")
		return MessageResult{}, fmt.Errorf("could not delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		newLogRequest(ctx, userID, "delete_user", "failed")
		return MessageResult{}, ErrUserNotFound
	}

	newLogRequest(ctx, userID, "delete_user", "success")
	return MessageResult{Message: "user deleted successfully"}, nil
}

func (s *AuthService) Heartbeat(_ context.Context, instanceID string) MessageResult {
	return MessageResult{Message: "heartbeat received"}
}

func (s *AuthService) getUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func newLogRequest(ctx context.Context, username, action, status string) {
	payload := map[string]string{
		"timestamp": time.Now().Format(time.RFC3339),
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
