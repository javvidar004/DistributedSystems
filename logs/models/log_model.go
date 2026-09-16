package models

import "time"

type Log struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Timestamp string    `json:"timestamp"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
