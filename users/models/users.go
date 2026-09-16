package models

import "time"

type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"not null;unique" json:"email"`
	Name         string    `gorm:"not null" json:"name"`
	LastName     string    `gorm:"not null" json:"last_name"`
	WorkPosition string    `gorm:"not null" json:"work_position"`
	Salary       float64   `gorm:"not null" json:"salary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
