package models

import "time"

type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	PasswordHash string    `gorm:"not null" json:"password"`
	Email        string    `gorm:"not null;unique" json:"email"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	Role         string    `gorm:"not null" json:"role"`
	Group        string    `gorm:"not null" json:"group"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
