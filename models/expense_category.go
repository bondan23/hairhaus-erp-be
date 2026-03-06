package models

// ExpenseCategory represents a category for operational expenses.
type ExpenseCategory struct {
	BaseModel
	Name string `gorm:"not null" json:"name"`
	Code string `gorm:"uniqueIndex;not null" json:"code"`
}
