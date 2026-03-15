package models

// Customer represents a walk-in or loyalty customer.
type Customer struct {
	BaseModel
	Name              string  `gorm:"not null" json:"name"`
	Phone             string  `json:"phone"`
	Gender            *string `json:"gender"`
	LoyaltyUserID     *string `json:"loyalty_user_id"`
	LoyaltyOutletID   *string `json:"loyalty_outlet_id"`
	IsLoyaltyVerified bool    `gorm:"default:false" json:"is_loyalty_verified"`
}
