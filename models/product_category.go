package models

// ProductCategory represents a category for products/services.
type ProductCategory struct {
	BaseModel
	Name       string `gorm:"not null" json:"name"`
	IncomeType string `json:"income_type"`
}
