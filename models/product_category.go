package models

const (
	CategoryCodeHaircut = "HAIRCUT"
)

// ProductCategory represents a category for products/services.
type ProductCategory struct {
	BaseModel
	Name string `gorm:"not null" json:"name"`
	Code string `gorm:"uniqueIndex;not null" json:"code"`
}
