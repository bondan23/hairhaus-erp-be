package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	DrawerStatusOpen    = "OPEN"
	DrawerStatusClosing = "CLOSING"
	DrawerStatusClosed  = "CLOSED"
)

// CashDrawer represents a cash drawer session for a branch.
type CashDrawer struct {
	BaseModel
	BranchID uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	Branch   Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`

	OpenedAt time.Time  `gorm:"not null" json:"opened_at"`
	ClosedAt *time.Time `json:"closed_at"`

	OpeningAmount int64 `gorm:"not null;default:0" json:"opening_amount"`

	ExpectedCash int64 `gorm:"not null;default:0" json:"expected_cash"`
	CountedCash  int64 `gorm:"not null;default:0" json:"counted_cash"`
	Variance     int64 `gorm:"not null;default:0" json:"variance"`

	Status string `gorm:"not null;default:'OPEN'" json:"status"` // OPEN | CLOSING | CLOSED

	ClosingSnapshot json.RawMessage `gorm:"type:jsonb" json:"closing_snapshot,omitempty"`
}
