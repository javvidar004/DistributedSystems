package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"

	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

var (
	ErrRequiredFields    = errors.New("required fields are missing")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("email already exists")
)

type UsersService struct {
	db *gorm.DB
}

type CreateUserInput struct {
	Email        string
	Name         string
	LastName     string
	WorkPosition string
	Salary       float64
}

type UpdateUserInput struct {
	Email        string
	Name         string
	LastName     string
	WorkPosition string
	Salary       float64
}

type UserResult struct {
	Message string
	User    models.User
}

type UsersResult struct {
	Users []models.User
}

type MessageResult struct {
	Message string
}

func NewUsersService(db *gorm.DB) *UsersService {
	return &UsersService{db: db}
}

func (s *UsersService) CreateUser(ctx context.Context, input CreateUserInput) (UserResult, error) {
	if input.Email == "" {
		return UserResult{}, ErrRequiredFields
	}

	if _, err := s.getUserByEmail(ctx, input.Email); err == nil {
		return UserResult{}, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return UserResult{}, err
	}

	user := models.User{
		Email:        input.Email,
		Name:         input.Name,
		LastName:     input.LastName,
		WorkPosition: input.WorkPosition,
		Salary:       input.Salary,
	}

	if err := s.db.WithContext(ctx).Create(&user).Error; err != nil {
		return UserResult{}, fmt.Errorf("could not create user: %w", err)
	}

	return UserResult{Message: "user created successfully", User: user}, nil
}

func (s *UsersService) UpdateUser(ctx context.Context, input UpdateUserInput) (UserResult, error) {
	if input.Email == "" {
		newLogRequest(ctx, input.Email, "update_user", "failed")
		return UserResult{}, ErrRequiredFields
	}

	user, err := s.getUserByEmail(ctx, input.Email)
	if err != nil {
		newLogRequest(ctx, input.Email, "update_user", "failed")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserResult{}, ErrUserNotFound
		}
		return UserResult{}, err
	}

	user.Name = input.Name
	user.LastName = input.LastName
	user.WorkPosition = input.WorkPosition
	user.Salary = input.Salary

	if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
		newLogRequest(ctx, input.Email, "update_user", "failed")
		return UserResult{}, fmt.Errorf("could not update user: %w", err)
	}

	newLogRequest(ctx, input.Email, "update_user", "success")
	return UserResult{Message: "user updated successfully", User: user}, nil
}

func (s *UsersService) DeleteUser(ctx context.Context, userID string) (MessageResult, error) {
	if userID == "" {
		return MessageResult{}, ErrRequiredFields
	}

	result := s.db.WithContext(ctx).Where("id = ?", userID).Delete(&models.User{})
	if result.Error != nil {
		return MessageResult{}, fmt.Errorf("could not delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return MessageResult{}, ErrUserNotFound
	}

	return MessageResult{Message: "user deleted successfully"}, nil
}

func (s *UsersService) GetUser(ctx context.Context, userID string) (UserResult, error) {
	if userID == "" {
		return UserResult{}, ErrRequiredFields
	}

	user, err := s.getUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return UserResult{}, ErrUserNotFound
		}
		return UserResult{}, err
	}

	return UserResult{Message: "user found successfully", User: user}, nil
}

func (s *UsersService) GetUsers(ctx context.Context) (UsersResult, error) {
	var users []models.User
	result := s.db.WithContext(ctx).Select("id", "email", "name", "last_name", "work_position", "salary", "created_at", "updated_at").Find(&users)
	if result.Error != nil {
		return UsersResult{}, fmt.Errorf("could not get users: %w", result.Error)
	}
	return UsersResult{Users: users}, nil
}

func (s *UsersService) Heartbeat(_ context.Context) MessageResult {
	return MessageResult{Message: "heartbeat received"}
}

func (s *UsersService) getUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UsersService) getUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func newLogRequest(ctx context.Context, username, action, status string) {
	middlewareHost := os.Getenv("MIDDLEWARE_HOST")
	if middlewareHost == "" {
		middlewareHost = "localhost"
	}

	dialCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		middlewareHost+":50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Printf("error dialing logs gRPC service: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewLogsServiceClient(conn)
	_, err = client.NewLog(ctx, &pb.NewLogRequest{
		Timestamp: time.Now().Format(time.RFC3339),
		Username:  username,
		Action:    action,
		Status:    status,
	})
	if err != nil {
		log.Printf("error sending log request over gRPC: %v", err)
		return
	}

	log.Printf("log request sent over gRPC: %s", action)
}
