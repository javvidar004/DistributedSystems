package models

import "time"

type Booking struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	WorkspaceID uint      `gorm:"not null;index" json:"workspace_id"`
	BookingDate time.Time `gorm:"type:date;not null;index" json:"booking_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Workspace   Workspace `gorm:"foreignKey:WorkspaceID" json:"workspace,omitempty"`
}
