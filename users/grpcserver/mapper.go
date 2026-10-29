package grpcserver

import (
	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

func userToProto(user models.User) *pb.User {
	return &pb.User{
		Id:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		LastName:     user.LastName,
		WorkPosition: user.WorkPosition,
		Salary:       user.Salary,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
