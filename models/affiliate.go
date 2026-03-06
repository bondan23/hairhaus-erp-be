package models

const (
	CommissionTypePercentage = "PERCENTAGE"
	CommissionTypeFixed      = "FIXED"
)

// Affiliate represents a referral partner linked to loyalty system.
type Affiliate struct {
	BaseModel
	LoyaltyMemberID string `gorm:"uniqueIndex;not null" json:"loyalty_member_id"`
	AffiliateCode   string `gorm:"uniqueIndex;not null" json:"affiliate_code"`
	Name            string `gorm:"not null" json:"name"`

	CommissionType       string  `gorm:"not null" json:"commission_type"` // PERCENTAGE | FIXED
	CommissionPercentage float64 `gorm:"default:0" json:"commission_percentage"`
	CommissionFixed      int64   `gorm:"default:0" json:"commission_fixed"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}
