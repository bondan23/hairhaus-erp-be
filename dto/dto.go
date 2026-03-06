package dto

import "github.com/google/uuid"

// ===== Auth =====
type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Pin         string `json:"pin" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ===== Branch =====
type CreateBranchRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type UpdateBranchRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

// ===== User =====
type CreateUserRequest struct {
	EmployeeID  string    `json:"employee_id" binding:"required"`
	OutletID    string    `json:"outlet_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Pin         string    `json:"pin" binding:"required,min=4"`
	Role        string    `json:"role" binding:"required,oneof=ADMIN MANAGER CASHIER"`
	BranchID    uuid.UUID `json:"branch_id" binding:"required"`
}

type UpdateUserRequest struct {
	EmployeeID  *string    `json:"employee_id"`
	OutletID    *string    `json:"outlet_id"`
	Name        *string    `json:"name"`
	PhoneNumber *string    `json:"phone_number"`
	Pin         *string    `json:"pin"`
	Role        *string    `json:"role"`
	BranchID    *uuid.UUID `json:"branch_id"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	EmployeeID  string    `json:"employee_id"`
	OutletID    string    `json:"outlet_id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	BranchID    uuid.UUID `json:"branch_id"`
}

// ===== Product Category =====
type CreateProductCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdateProductCategoryRequest struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
}

// ===== Product =====
type CreateProductRequest struct {
	Name        string    `json:"name" binding:"required"`
	ProductType string    `json:"product_type" binding:"required,oneof=SERVICE RETAIL"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	BasePrice   int64     `json:"base_price" binding:"min=0"`
	CostPrice   int64     `json:"cost_price" binding:"min=0"`
}

type UpdateProductRequest struct {
	Name        *string    `json:"name"`
	ProductType *string    `json:"product_type"`
	CategoryID  *uuid.UUID `json:"category_id"`
	BasePrice   *int64     `json:"base_price"`
	CostPrice   *int64     `json:"cost_price"`
	IsActive    *bool      `json:"is_active"`
}

// ===== Stylist =====
type CreateStylistRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateStylistRequest struct {
	Name     *string `json:"name"`
	IsActive *bool   `json:"is_active"`
}

// ===== Branch Product =====
type CreateBranchProductRequest struct {
	BranchID      uuid.UUID `json:"branch_id" binding:"required"`
	ProductID     uuid.UUID `json:"product_id" binding:"required"`
	PriceOverride *int64    `json:"price_override"`
	Stock         int64     `json:"stock"`
}

type UpdateBranchProductRequest struct {
	PriceOverride *int64 `json:"price_override"`
}

// ===== Branch Stylist =====
type CreateBranchStylistRequest struct {
	BranchID             uuid.UUID `json:"branch_id" binding:"required"`
	StylistID            uuid.UUID `json:"stylist_id" binding:"required"`
	HaircutPriceOverride *int64    `json:"haircut_price_override"`
	CommissionPercentage *int      `json:"commission_percentage"`
}

type UpdateBranchStylistRequest struct {
	HaircutPriceOverride *int64 `json:"haircut_price_override"`
	CommissionPercentage *int   `json:"commission_percentage"`
}

// ===== Customer =====
type CreateCustomerRequest struct {
	Name              string `json:"name" binding:"required"`
	Phone             string `json:"phone"`
	LoyaltyExternalID string `json:"loyalty_external_id"`
}

type UpdateCustomerRequest struct {
	Name              *string `json:"name"`
	Phone             *string `json:"phone"`
	LoyaltyExternalID *string `json:"loyalty_external_id"`
}

// ===== Cash Drawer =====
type OpenDrawerRequest struct {
	BranchID      uuid.UUID `json:"branch_id" binding:"required"`
	OpeningAmount int64     `json:"opening_amount" binding:"min=0"`
}

type CloseDrawerRequest struct {
	CountedCash int64 `json:"counted_cash" binding:"min=0"`
}

// ===== Transaction / POS =====
type CheckoutRequest struct {
	BranchID       uuid.UUID          `json:"branch_id" binding:"required"`
	CustomerID     *uuid.UUID         `json:"customer_id"`
	AffiliateCode  string             `json:"affiliate_code"`
	Items          []CheckoutItemRequest `json:"items" binding:"required,min=1,dive"`
	Payments       []PaymentRequest   `json:"payments" binding:"required,min=1,dive"`
	DiscountAmount int64              `json:"discount_amount" binding:"min=0"`
	IdempotencyKey string             `json:"idempotency_key" binding:"required"`
}

type CheckoutItemRequest struct {
	ProductID uuid.UUID  `json:"product_id" binding:"required"`
	StylistID *uuid.UUID `json:"stylist_id"`
	Quantity  int64      `json:"quantity" binding:"required,min=1"`
}

type PaymentRequest struct {
	Method      string `json:"method" binding:"required,oneof=CASH QRIS DEBIT CREDIT"`
	Amount      int64  `json:"amount" binding:"required,min=1"`
	ReferenceNo string `json:"reference_no"`
}

// ===== Transaction Edit =====
type EditTransactionRequest struct {
	DiscountAmount *int64  `json:"discount_amount"`
	EditReason     string  `json:"edit_reason" binding:"required"`
}

type EditPaymentRequest struct {
	Payments   []PaymentRequest `json:"payments" binding:"required,min=1,dive"`
	EditReason string           `json:"edit_reason" binding:"required"`
}

// ===== Inventory =====
type RestockRequest struct {
	BranchID  uuid.UUID `json:"branch_id" binding:"required"`
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int64     `json:"quantity" binding:"required,min=1"`
	Note      string    `json:"note"`
}

type AdjustStockRequest struct {
	BranchID  uuid.UUID `json:"branch_id" binding:"required"`
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Change    int64     `json:"change" binding:"required"`
	Note      string    `json:"note" binding:"required"`
}

// ===== Affiliate =====
type CreateAffiliateRequest struct {
	LoyaltyMemberID      string  `json:"loyalty_member_id" binding:"required"`
	AffiliateCode        string  `json:"affiliate_code" binding:"required"`
	Name                 string  `json:"name" binding:"required"`
	CommissionType       string  `json:"commission_type" binding:"required,oneof=PERCENTAGE FIXED"`
	CommissionPercentage float64 `json:"commission_percentage"`
	CommissionFixed      int64   `json:"commission_fixed"`
}

type UpdateAffiliateRequest struct {
	Name                 *string  `json:"name"`
	CommissionType       *string  `json:"commission_type"`
	CommissionPercentage *float64 `json:"commission_percentage"`
	CommissionFixed      *int64   `json:"commission_fixed"`
	IsActive             *bool    `json:"is_active"`
}

// ===== Salary =====
type GenerateSalaryRequest struct {
	BranchID  uuid.UUID `json:"branch_id" binding:"required"`
	StylistID uuid.UUID `json:"stylist_id" binding:"required"`
	Month     int       `json:"month" binding:"required,min=1,max=12"`
	Year      int       `json:"year" binding:"required,min=2020"`
}

// ===== Expense Category =====
type CreateExpenseCategoryRequest struct {
	Name string `json:"name" binding:"required"`
	Code string `json:"code" binding:"required"`
}

type UpdateExpenseCategoryRequest struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
}

// ===== Expense =====
type CreateExpenseRequest struct {
	BranchID    uuid.UUID `json:"branch_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Description string    `json:"description"`
	Amount      int64     `json:"amount" binding:"required,min=1"`
	ExpenseDate string    `json:"expense_date" binding:"required"` // YYYY-MM-DD
}

type UpdateExpenseRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Description *string    `json:"description"`
	Amount      *int64     `json:"amount"`
	ExpenseDate *string    `json:"expense_date"`
}

// ===== Report =====
type ReportFilter struct {
	BranchID  uuid.UUID `form:"branch_id"`
	StartDate string    `form:"start_date"` // YYYY-MM-DD
	EndDate   string    `form:"end_date"`   // YYYY-MM-DD
}

type FinancialReport struct {
	Revenue     int64 `json:"revenue"`
	COGS        COGSBreakdown `json:"cogs"`
	GrossProfit int64 `json:"gross_profit"`
	OPEX        int64 `json:"opex"`
	NetProfit   int64 `json:"net_profit"`
}

type COGSBreakdown struct {
	Commission          int64 `json:"commission"`
	AffiliateCommission int64 `json:"affiliate_commission"`
	RetailCost          int64 `json:"retail_cost"`
	Discount            int64 `json:"discount"`
	Total               int64 `json:"total"`
}
