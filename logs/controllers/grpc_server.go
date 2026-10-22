package controllers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

type LogGRPCServer struct {
	pb.UnimplementedLogServiceServer
	db *gorm.DB
}

func NewLogGRPCServer(db *gorm.DB) *LogGRPCServer {
	return &LogGRPCServer{db: db}
}

func (s *LogGRPCServer) CreateLog(ctx context.Context, req *pb.CreateLogRequest) (*pb.CreateLogResponse, error) {
	_ = ctx
	if req.GetAction() == "" {
		return nil, status.Error(codes.InvalidArgument, "action is required")
	}

	entry := models.Log{
		ID:        uuid.New().String(),
		Timestamp: time.Now().Format(time.RFC3339),
		Username:  req.GetUserId(),
		Action:    req.GetAction(),
		Status:    req.GetDetails(),
	}

	if err := s.db.Create(&entry).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not create log")
	}

	return &pb.CreateLogResponse{
		Success: true,
		LogId:   entry.ID,
	}, nil
}

func (s *LogGRPCServer) GetLogs(ctx context.Context, req *pb.GetLogsRequest) (*pb.GetLogsResponse, error) {
	_ = ctx
	query := s.db.Order("created_at desc")
	if req.GetLimit() > 0 {
		query = query.Limit(int(req.GetLimit()))
	}

	var logs []models.Log
	if err := query.Find(&logs).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not fetch logs")
	}

	resp := &pb.GetLogsResponse{Logs: make([]*pb.LogEntry, 0, len(logs))}
	for _, l := range logs {
		resp.Logs = append(resp.Logs, &pb.LogEntry{
			Id:          l.ID,
			ServiceName: "",
			Action:      l.Action,
			UserId:      l.Username,
			Timestamp:   l.Timestamp,
			Details:     l.Status,
		})
	}

	return resp, nil
}
