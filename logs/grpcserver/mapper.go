package grpcserver

import (
	"proyecto.com/m/models"
	"proyecto.com/m/pb"
)

func logToProto(logEntry models.Log) *pb.Log {
	return &pb.Log{
		Id:        logEntry.ID,
		Timestamp: logEntry.Timestamp,
		Username:  logEntry.Username,
		Action:    logEntry.Action,
		Status:    logEntry.Status,
		CreatedAt: logEntry.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
