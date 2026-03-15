package models

// Branch represents a physical barbershop location.
type Branch struct {
	BaseModel
	Name     string `gorm:"not null" json:"name"`
	Code     string `gorm:"uniqueIndex;not null" json:"code"`
	OutletID string `gorm:"not null" json:"outlet_id"` // HAIRHAUS Loyalty Outlet ID
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
