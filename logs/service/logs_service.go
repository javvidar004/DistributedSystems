package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"proyecto.com/m/models"
)

var (
	ErrInvalidPayload = errors.New("invalid log payload")
)

type LogsService struct {
	db *gorm.DB
}

type NewLogInput struct {
	Timestamp string
	Username  string
	Action    string
	Status    string
}

type LogResult struct {
	Message string
	Log     models.Log
}

type LogsResult struct {
	Logs []models.Log
}

type MessageResult struct {
	Message string
}

func NewLogsService(db *gorm.DB) *LogsService {
	return &LogsService{db: db}
}

func (s *LogsService) GetLogs(ctx context.Context) (LogsResult, error) {
	var logs []models.Log
	if err := s.db.WithContext(ctx).Find(&logs).Error; err != nil {
		return LogsResult{}, fmt.Errorf("could not fetch logs: %w", err)
	}
	return LogsResult{Logs: logs}, nil
}

func (s *LogsService) NewLog(ctx context.Context, input NewLogInput) (LogResult, error) {
	logEntry := models.Log{
		ID:        uuid.New().String(),
		Timestamp: input.Timestamp,
		Username:  input.Username,
		Action:    input.Action,
		Status:    input.Status,
	}

	if err := s.db.WithContext(ctx).Create(&logEntry).Error; err != nil {
		return LogResult{}, fmt.Errorf("could not create log: %w", err)
	}

	return LogResult{Message: "log created successfully", Log: logEntry}, nil
}

func (s *LogsService) Heartbeat(_ context.Context) MessageResult {
	return MessageResult{Message: "heartbeat received"}
}
