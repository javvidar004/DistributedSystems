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

type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserGRPCServer(db *gorm.DB) *UserGRPCServer {
	return &UserGRPCServer{db: db}
}

func (s *UserGRPCServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	_ = ctx
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if _, err := getUserByEmail(s.db, req.GetEmail()); err == nil {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	user := models.User{
		Email:        req.GetEmail(),
		Name:         req.GetUsername(),
		LastName:     "",
		WorkPosition: req.GetRole(),
		Salary:       0,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not create user")
	}

	return &pb.UserResponse{
		User:    modelToPBUser(user),
		Message: "user created successfully",
	}, nil
}

func (s *UserGRPCServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	_ = ctx
	user, err := getUserByID(s.db, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "could not fetch user")
	}

	return &pb.UserResponse{
		User:    modelToPBUser(user),
		Message: "user found successfully",
	}, nil
}

func (s *UserGRPCServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	_ = ctx
	var users []models.User

	query := s.db.Model(&models.User{})

	limit := int(req.GetLimit())
	if limit > 0 {
		query = query.Limit(limit)
	}

	page := int(req.GetPage())
	if limit > 0 && page > 1 {
		query = query.Offset((page - 1) * limit)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not list users")
	}

	var total int64
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not count users")
	}

	pbUsers := make([]*pb.User, 0, len(users))
	for _, u := range users {
		pbUsers = append(pbUsers, modelToPBUser(u))
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: int32(total),
	}, nil
}

func (s *UserGRPCServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	_ = ctx
	user, err := getUserByID(s.db, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "could not fetch user")
	}

	if req.GetUsername() != "" {
		user.Name = req.GetUsername()
	}
	if req.GetEmail() != "" {
		user.Email = req.GetEmail()
	}
	if req.GetRole() != "" {
		user.WorkPosition = req.GetRole()
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not update user")
	}

	return &pb.UserResponse{
		User:    modelToPBUser(user),
		Message: "user updated successfully",
	}, nil
}

func (s *UserGRPCServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_ = ctx
	result := s.db.Where("id = ?", req.GetId()).Delete(&models.User{})
	if result.Error != nil {
		return nil, status.Error(codes.Internal, "could not delete user")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.DeleteUserResponse{Success: true, Message: "user deleted successfully"}, nil
}

func modelToPBUser(user models.User) *pb.User {
	return &pb.User{
		Id:       user.ID,
		Username: user.Name,
		Email:    user.Email,
		Role:     user.WorkPosition,
	}
}
