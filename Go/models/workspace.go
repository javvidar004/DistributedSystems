package models

type Workspace struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Number   string `gorm:"not null" json:"number"`
	Level    int    `gorm:"not null" json:"level"`
	OfficeID uint   `gorm:"not null;index" json:"office_id"`
	Office   Office `gorm:"foreignKey:OfficeID" json:"office,omitempty"`
}
