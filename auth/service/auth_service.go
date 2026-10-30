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
	"proyecto.com/m/pb"
)

var (
	ErrInvalidPayload      = errors.New("invalid payload")
	ErrRequiredFields      = errors.New("required fields are missing")
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("username already exists")
	ErrOldPasswordMismatch = errors.New("old password is incorrect")
	ErrReplicationFailed   = errors.New("failed to replicate user in users service")
)

type AuthService struct {
	db          *gorm.DB
	usersClient pb.UsersServiceClient
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

func NewAuthService(db *gorm.DB, usersClient pb.UsersServiceClient) *AuthService {
	return &AuthService{db: db, usersClient: usersClient}
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

	matched, needsMigration := verifyPassword(user.PasswordHash, password)
	if !matched {
		newLogRequest(ctx, email, "login", "failed")
		return LoginResult{}, ErrInvalidCredentials
	}

	if needsMigration {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			newLogRequest(ctx, email, "login", "failed")
			return LoginResult{}, fmt.Errorf("could not migrate legacy password hash: %w", err)
		}

		if err := s.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Update("password_hash", hashedPassword).Error; err != nil {
			newLogRequest(ctx, email, "login", "failed")
			return LoginResult{}, fmt.Errorf("could not persist migrated password hash: %w", err)
		}

		user.PasswordHash = hashedPassword
	}

	newLogRequest(ctx, email, "login", "success")
	return LoginResult{Message: "login successful", User: user}, nil
}

func (s *AuthService) RegisterAndReplicate(ctx context.Context, input RegisterInput) (RegisterResult, error) {
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

	hash, err := hashPassword(input.Password)
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
		PasswordHash: hash,
		IsActive:     true,
		Role:         "user",
		Group:        "default",
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, fmt.Errorf("could not create user: %w", err)
	}

	if err := s.replicateUserToUsersService(ctx, user.ID, input); err != nil {
		// Best-effort compensation to avoid inconsistent state between auth and users DBs.
		if delErr := s.db.WithContext(ctx).Delete(&models.User{}, "id = ?", user.ID).Error; delErr != nil {
			log.Printf("failed to rollback auth user %s after users replication error: %v", user.ID, delErr)
		}
		newLogRequest(ctx, input.Email, "register", "failed")
		return RegisterResult{}, err
	}

	newLogRequest(ctx, input.Email, "register", "success")
	return RegisterResult{Message: "user created successfully", User: user}, nil
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (RegisterResult, error) {
	return s.RegisterAndReplicate(ctx, input)
}

func (s *AuthService) replicateUserToUsersService(ctx context.Context, userID string, input RegisterInput) error {
	if s.usersClient == nil {
		return ErrReplicationFailed
	}

	grpcCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.usersClient.CreateUser(grpcCtx, &pb.CreateUserRequest{
		Id:           userID,
		Email:        input.Email,
		Name:         input.Name,
		LastName:     input.LastName,
		WorkPosition: input.WorkPosition,
		Salary:       input.Salary,
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrReplicationFailed, err)
	}

	return nil
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

	matched, _ := verifyPassword(user.PasswordHash, oldPassword)
	if !matched {
		newLogRequest(ctx, email, "update_password", "failed")
		return MessageResult{}, ErrOldPasswordMismatch
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		newLogRequest(ctx, email, "update_password", "failed")
		return MessageResult{}, fmt.Errorf("could not hash new password: %w", err)
	}

	user.PasswordHash = hashedPassword
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

func hashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(hash)), nil
}

func verifyPassword(storedPasswordHash, plainPassword string) (bool, bool) {
	if isBcryptHash(storedPasswordHash) {
		err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(plainPassword))
		return err == nil, false
	}

	return storedPasswordHash == plainPassword, true
}

func isBcryptHash(value string) bool {
	return strings.HasPrefix(value, "$2a$") || strings.HasPrefix(value, "$2b$") || strings.HasPrefix(value, "$2y$")
}
