package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	LastName  string    `gorm:"not null" json:"last_name"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Email     string    `gorm:"not null;unique" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
