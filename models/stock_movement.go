package models

import "github.com/google/uuid"

const (
	StockMovementTypeSale       = "SALE"
	StockMovementTypeRestock    = "RESTOCK"
	StockMovementTypeAdjustment = "ADJUSTMENT"
)

// StockMovement represents an inventory ledger entry.
type StockMovement struct {
	BaseModel
	BranchID  uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`

	Change int64  `gorm:"not null" json:"change"`
	Type   string `gorm:"not null" json:"type"` // SALE | RESTOCK | ADJUSTMENT

	ReferenceID uuid.UUID `gorm:"type:uuid" json:"reference_id"`
	Note        string    `json:"note"`
}
