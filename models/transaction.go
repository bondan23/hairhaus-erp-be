package models

import "github.com/google/uuid"

const (
	TransactionStatusDraft     = "DRAFT"
	TransactionStatusCompleted = "COMPLETED"
	TransactionStatusVoided    = "VOIDED"
)

// Transaction represents a POS transaction.
type Transaction struct {
	BaseModel
	InvoiceNo string `gorm:"uniqueIndex;not null" json:"invoice_no"`

	BranchID uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	Branch   Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`

	CustomerID *uuid.UUID `gorm:"type:uuid" json:"customer_id"`
	Customer   *Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`

	CustomerName *string `json:"customer_name"`

	AffiliateID *uuid.UUID `gorm:"type:uuid" json:"affiliate_id"`
	Affiliate   *Affiliate `gorm:"foreignKey:AffiliateID" json:"affiliate,omitempty"`

	SubtotalAmount int64 `gorm:"not null;default:0" json:"subtotal_amount"`
	DiscountAmount int64 `gorm:"not null;default:0" json:"discount_amount"`
	TotalAmount    int64 `gorm:"not null;default:0" json:"total_amount"`

	AffiliateCommissionAmountSnapshot int64 `gorm:"not null;default:0" json:"affiliate_commission_amount_snapshot"`

	Status string `gorm:"not null;default:'DRAFT'" json:"status"` // DRAFT | COMPLETED | VOIDED

	CashDrawerID uuid.UUID `gorm:"type:uuid;not null" json:"cash_drawer_id"`
	// CashDrawer   CashDrawer `gorm:"foreignKey:CashDrawerID" json:"cash_drawer,omitempty"`

	IdempotencyKey string `gorm:"uniqueIndex" json:"idempotency_key,omitempty"`

	EditedByID *uuid.UUID `gorm:"type:uuid" json:"edited_by_id"`
	EditReason string     `json:"edit_reason"`

	CreatedByID uuid.UUID `gorm:"type:uuid;not null" json:"created_by_id"`

	Items    []TransactionItem `gorm:"foreignKey:TransactionID" json:"items,omitempty"`
	Payments []Payment         `gorm:"foreignKey:TransactionID" json:"payments,omitempty"`
}
