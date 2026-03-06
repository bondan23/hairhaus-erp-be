package models

import "github.com/google/uuid"

const (
	SalaryStatusGenerated = "GENERATED"
	SalaryStatusPaid      = "PAID"
)

// SalaryRecord represents a monthly salary calculation for a stylist at a branch.
type SalaryRecord struct {
	BaseModel
	StylistID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_salary_unique" json:"stylist_id"`
	BranchID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_salary_unique" json:"branch_id"`

	Month int `gorm:"not null;uniqueIndex:idx_salary_unique" json:"month"`
	Year  int `gorm:"not null;uniqueIndex:idx_salary_unique" json:"year"`

	TotalSales      int64 `gorm:"not null;default:0" json:"total_sales"`
	TotalCommission int64 `gorm:"not null;default:0" json:"total_commission"`

	Status string `gorm:"not null;default:'GENERATED'" json:"status"` // GENERATED | PAID

	Stylist Stylist `gorm:"foreignKey:StylistID" json:"stylist,omitempty"`
	Branch  Branch  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
