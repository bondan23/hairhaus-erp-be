package models

// Customer represents a walk-in or loyalty customer.
type Customer struct {
	BaseModel
	Name              string `gorm:"not null" json:"name"`
	Phone             string `json:"phone"`
	LoyaltyExternalID string `json:"loyalty_external_id"`
}
