package grpcserver

import (
	"context"
	"errors"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"proyecto.com/m/pb"
	"proyecto.com/m/service"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authSvc *service.AuthService
}

func NewAuthServer(authSvc *service.AuthService) *AuthServer {
	return &AuthServer{authSvc: authSvc}
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	res, err := s.authSvc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return &pb.AuthResponse{
		Message: res.Message,
		User:    userToProto(res.User),
	}, nil
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	res, err := s.authSvc.Register(ctx, service.RegisterInput{
		Name:         req.GetName(),
		LastName:     req.GetLastName(),
		Email:        req.GetEmail(),
		WorkPosition: req.GetWorkPosition(),
		Salary:       req.GetSalary(),
		Password:     req.GetPassword(),
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return &pb.AuthResponse{
		Message: res.Message,
		User:    userToProto(res.User),
	}, nil
}

func (s *AuthServer) UpdatePassword(ctx context.Context, req *pb.UpdatePasswordRequest) (*pb.MessageResponse, error) {
	res, err := s.authSvc.UpdatePassword(ctx, req.GetEmail(), req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return &pb.MessageResponse{Message: res.Message}, nil
}

func (s *AuthServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.MessageResponse, error) {
	res, err := s.authSvc.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, mapGRPCError(err)
	}

	return &pb.MessageResponse{Message: res.Message}, nil
}

func (s *AuthServer) Heartbeat(ctx context.Context, _ *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	res := s.authSvc.Heartbeat(ctx, os.Getenv("INSTANCE_ID"))
	return &pb.HeartbeatResponse{
		Message:    res.Message,
		InstanceId: os.Getenv("INSTANCE_ID"),
	}, nil
}

func mapGRPCError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidPayload), errors.Is(err, service.ErrRequiredFields):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials), errors.Is(err, service.ErrOldPasswordMismatch):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, service.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
