package models

// Stylist represents a hair maker / barber.
type Stylist struct {
	BaseModel
	Name     string `gorm:"not null" json:"name"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
