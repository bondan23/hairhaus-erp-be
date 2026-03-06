package models

import "github.com/google/uuid"

// BranchProduct represents branch-level price override and inventory for a product.
type BranchProduct struct {
	BaseModel
	BranchID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_branch_product" json:"branch_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_branch_product" json:"product_id"`

	PriceOverride *int64 `json:"price_override"`

	Stock int64 `gorm:"not null;default:0" json:"stock"`

	Branch  Branch  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
