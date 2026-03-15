package models

import "github.com/google/uuid"

const (
	RoleAdmin   = "ADMIN"
	RoleManager = "MANAGER"
	RoleCashier = "CASHIER"
)

// User represents a system user (admin, manager, or cashier).
type User struct {
	BaseModel
	LoyaltyEmployeeID string `gorm:"uniqueIndex;not null" json:"loyalty_employee_id"` // HAIRHAUS Loyalty ID, manually input
	Name              string `gorm:"not null" json:"name"`
	PhoneNumber       string `gorm:"uniqueIndex;not null" json:"phone_number"`
	Pin               string `gorm:"not null" json:"-"` // hashed PIN

	Role string `gorm:"not null" json:"role"`

	BranchID uuid.UUID `gorm:"type:uuid;not null" json:"branch_id"`
	Branch   Branch    `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
