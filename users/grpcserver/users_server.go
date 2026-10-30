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

type UsersServer struct {
	pb.UnimplementedUsersServiceServer
	usersSvc *service.UsersService
}

func NewUsersServer(usersSvc *service.UsersService) *UsersServer {
	return &UsersServer{usersSvc: usersSvc}
}

func (s *UsersServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	res, err := s.usersSvc.CreateUser(ctx, service.CreateUserInput{
		ID:           req.GetId(),
		Email:        req.GetEmail(),
		Name:         req.GetName(),
		LastName:     req.GetLastName(),
		WorkPosition: req.GetWorkPosition(),
		Salary:       req.GetSalary(),
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return &pb.UserResponse{Message: res.Message, User: userToProto(res.User)}, nil
}

func (s *UsersServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	res, err := s.usersSvc.UpdateUser(ctx, service.UpdateUserInput{
		Email:        req.GetEmail(),
		Name:         req.GetName(),
		LastName:     req.GetLastName(),
		WorkPosition: req.GetWorkPosition(),
		Salary:       req.GetSalary(),
	})
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return &pb.UserResponse{Message: res.Message, User: userToProto(res.User)}, nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.MessageResponse, error) {
	res, err := s.usersSvc.DeleteUser(ctx, req.GetUserId())
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return &pb.MessageResponse{Message: res.Message}, nil
}

func (s *UsersServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	res, err := s.usersSvc.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, mapGRPCError(err)
	}
	return &pb.UserResponse{Message: res.Message, User: userToProto(res.User)}, nil
}

func (s *UsersServer) GetUsers(ctx context.Context, _ *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	res, err := s.usersSvc.GetUsers(ctx)
	if err != nil {
		return nil, mapGRPCError(err)
	}

	users := make([]*pb.User, 0, len(res.Users))
	for _, u := range res.Users {
		users = append(users, userToProto(u))
	}

	return &pb.GetUsersResponse{Users: users}, nil
}

func (s *UsersServer) Heartbeat(ctx context.Context, _ *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	res := s.usersSvc.Heartbeat(ctx)
	return &pb.HeartbeatResponse{Message: res.Message, InstanceId: os.Getenv("INSTANCE_ID")}, nil
}

func mapGRPCError(err error) error {
	switch {
	case errors.Is(err, service.ErrRequiredFields):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
