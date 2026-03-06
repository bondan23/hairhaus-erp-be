package models

import "github.com/google/uuid"

const (
	IncomeTypeHaircut   = "HAIRCUT"
	IncomeTypeTreatment = "TREATMENT"
	IncomeTypeProduct   = "PRODUCT"
)

// TransactionItem represents a line item in a transaction with snapshot fields.
type TransactionItem struct {
	BaseModel
	TransactionID uuid.UUID   `gorm:"type:uuid;not null" json:"transaction_id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionID" json:"-"`

	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`

	StylistID *uuid.UUID `gorm:"type:uuid" json:"stylist_id"`
	Stylist   *Stylist   `gorm:"foreignKey:StylistID" json:"stylist,omitempty"`

	// Snapshot fields for financial integrity
	ProductNameSnapshot  string `gorm:"not null" json:"product_name_snapshot"`
	ProductTypeSnapshot  string `gorm:"not null" json:"product_type_snapshot"`
	CategoryNameSnapshot string `gorm:"not null" json:"category_name_snapshot"`
	IncomeTypeSnapshot   string `gorm:"not null" json:"income_type_snapshot"` // HAIRCUT | TREATMENT | PRODUCT
	StylistNameSnapshot  string `json:"stylist_name_snapshot"`

	PriceSnapshot int64 `gorm:"not null" json:"price_snapshot"`
	Quantity      int64 `gorm:"not null;default:1" json:"quantity"`

	GrossSubtotal int64 `gorm:"not null" json:"gross_subtotal"`

	ItemDiscount int64 `gorm:"not null;default:0" json:"item_discount"`
	NetSubtotal  int64 `gorm:"not null" json:"net_subtotal"`

	CommissionAmountSnapshot int64 `gorm:"not null;default:0" json:"commission_amount_snapshot"`
	CostPriceSnapshot        int64 `gorm:"not null;default:0" json:"cost_price_snapshot"`
}
