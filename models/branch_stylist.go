package models

import "github.com/google/uuid"

// BranchStylist represents a stylist assignment to a branch with optional haircut price override.
type BranchStylist struct {
	BaseModel
	BranchID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_branch_stylist" json:"branch_id"`
	StylistID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_branch_stylist" json:"stylist_id"`

	HaircutPriceOverride *int64 `json:"haircut_price_override"`
	CommissionPercentage *int   `json:"commission_percentage"`

	Branch  Branch  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Stylist Stylist `gorm:"foreignKey:StylistID" json:"stylist,omitempty"`
}
