package models

type Office struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
