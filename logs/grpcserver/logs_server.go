package grpcserver

import (
	"context"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"proyecto.com/m/pb"
	"proyecto.com/m/service"
)

type LogsServer struct {
	pb.UnimplementedLogsServiceServer
	logsSvc *service.LogsService
}

func NewLogsServer(logsSvc *service.LogsService) *LogsServer {
	return &LogsServer{logsSvc: logsSvc}
}

func (s *LogsServer) GetLogs(ctx context.Context, _ *pb.GetLogsRequest) (*pb.GetLogsResponse, error) {
	res, err := s.logsSvc.GetLogs(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	logs := make([]*pb.Log, 0, len(res.Logs))
	for _, l := range res.Logs {
		logs = append(logs, logToProto(l))
	}

	return &pb.GetLogsResponse{Logs: logs}, nil
}

func (s *LogsServer) NewLog(ctx context.Context, req *pb.NewLogRequest) (*pb.LogResponse, error) {
	res, err := s.logsSvc.NewLog(ctx, service.NewLogInput{
		Timestamp: req.GetTimestamp(),
		Username:  req.GetUsername(),
		Action:    req.GetAction(),
		Status:    req.GetStatus(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.LogResponse{Message: res.Message, Log: logToProto(res.Log)}, nil
}

func (s *LogsServer) Heartbeat(ctx context.Context, _ *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	res := s.logsSvc.Heartbeat(ctx)
	return &pb.HeartbeatResponse{Message: res.Message, InstanceId: os.Getenv("INSTANCE_ID")}, nil
}
