package controllers

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

type AuthGRPCServer struct {
	pb.UnimplementedAuthServiceServer
	db *gorm.DB
}

func NewAuthGRPCServer(db *gorm.DB) *AuthGRPCServer {
	return &AuthGRPCServer{db: db}
}

func (s *AuthGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_ = ctx
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	var user models.User
	if err := s.db.Where("email = ?", req.GetUsername()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid username or password")
		}
		return nil, status.Error(codes.Internal, "failed to query user")
	}

	if user.PasswordHash != req.GetPassword() {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	return &pb.LoginResponse{
		Token:   user.ID,
		UserId:  user.ID,
		Role:    user.Role,
		Message: "login successful",
	}, nil
}

func (s *AuthGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	_ = ctx
	if req.GetToken() == "" {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	var user models.User
	if err := s.db.Where("id = ?", req.GetToken()).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.ValidateTokenResponse{Valid: false}, nil
		}
		return nil, status.Error(codes.Internal, "failed to validate token")
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: user.ID,
		Role:   user.Role,
	}, nil
}
