package models

import "github.com/google/uuid"

const (
	PaymentMethodCash   = "CASH"
	PaymentMethodQRIS   = "QRIS"
	PaymentMethodDebit  = "DEBIT"
	PaymentMethodCredit = "CREDIT"
)

// Payment represents a payment towards a transaction. Supports split payments.
type Payment struct {
	BaseModel
	TransactionID uuid.UUID   `gorm:"type:uuid;not null" json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"-"`

	Method string `gorm:"not null" json:"method"` // CASH | QRIS | DEBIT | CREDIT
	Amount int64  `gorm:"not null" json:"amount"`

	ReferenceNo string `json:"reference_no"`
}
