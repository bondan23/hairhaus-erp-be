package models

// Branch represents a physical barbershop location.
type Branch struct {
	BaseModel
	Name            string `gorm:"not null" json:"name"`
	Code            string `gorm:"uniqueIndex;not null" json:"code"`
	LoyaltyOutletID string `gorm:"not null" json:"loyalty_outlet_id"` // HAIRHAUS Loyalty Outlet ID
	Address         string `json:"address"`
	Phone           string `json:"phone"`
	IsActive        bool   `gorm:"default:true" json:"is_active"`
}
