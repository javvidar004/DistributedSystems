package grpcserver

import (
	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

func userToProto(user models.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Email:     user.Email,
		IsActive:  user.IsActive,
		Role:      user.Role,
		Group:     user.Group,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
