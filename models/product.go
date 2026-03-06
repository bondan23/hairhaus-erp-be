package models

import "github.com/google/uuid"

const (
	ProductTypeService = "SERVICE"
	ProductTypeRetail  = "RETAIL"
)

// Product represents a service or retail item.
type Product struct {
	BaseModel
	Name string `gorm:"not null" json:"name"`

	ProductType string `gorm:"not null" json:"product_type"` // SERVICE | RETAIL

	CategoryID uuid.UUID       `gorm:"type:uuid;not null" json:"category_id"`
	Category   ProductCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	BasePrice int64 `gorm:"not null;default:0" json:"base_price"`
	CostPrice int64 `gorm:"not null;default:0" json:"cost_price"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}
