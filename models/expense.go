package models

import (
	"time"

	"github.com/google/uuid"
)

// Expense represents an operational cost entry.
type Expense struct {
	BaseModel
	BranchID uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	Branch   Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`

	CategoryID uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
	Category   ExpenseCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	Description string    `json:"description"`
	Amount      int64     `gorm:"not null" json:"amount"`
	ExpenseDate time.Time `gorm:"not null" json:"expense_date"`
}
