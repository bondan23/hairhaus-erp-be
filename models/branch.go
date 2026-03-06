package models

// Branch represents a physical barbershop location.
type Branch struct {
	BaseModel
	Name     string `gorm:"not null" json:"name"`
	Code     string `gorm:"uniqueIndex;not null" json:"code"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
