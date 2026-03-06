package models

import "github.com/google/uuid"

const (
	AffiliateCommissionStatusPending = "PENDING"
	AffiliateCommissionStatusPaid    = "PAID"
)

// AffiliateCommission tracks commission earned by an affiliate from a transaction.
type AffiliateCommission struct {
	BaseModel
	AffiliateID   uuid.UUID   `gorm:"type:uuid;not null" json:"affiliate_id"`
	TransactionID uuid.UUID   `gorm:"type:uuid;not null" json:"transaction_id"`
	Affiliate     Affiliate   `gorm:"foreignKey:AffiliateID" json:"affiliate,omitempty"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"transaction,omitempty"`

	CommissionAmount int64 `gorm:"not null" json:"commission_amount"`

	Status string `gorm:"not null;default:'PENDING'" json:"status"` // PENDING | PAID
}
