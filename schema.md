# **Schema & ERD**

# **1\. BASE MODEL**

Common fields used in most tables.

type BaseModel struct {  
	ID        uuid.UUID      \`gorm:"type:uuid;default:gen\_random\_uuid();primaryKey"\`  
	CreatedAt time.Time  
	UpdatedAt time.Time  
	DeletedAt gorm.DeletedAt \`gorm:"index"\`  
}  
---

# **2\. BRANCH**

type Branch struct {  
	BaseModel

	Name    string  
	Code    string \`gorm:"uniqueIndex"\`  
	Address string  
	Phone   string  
	IsActive bool  
}  
---

# **3\. USER**

type User struct {  
	BaseModel

	Name     string  
	Email    string \`gorm:"uniqueIndex"\`  
	Password string

	Role string

	BranchID uuid.UUID  
	Branch   Branch  
}

Roles:

ADMIN  
MANAGER  
CASHIER  
---

# **4\. PRODUCT CATEGORY**

type ProductCategory struct {  
	BaseModel

	Name string  
	Code string \`gorm:"uniqueIndex"\`  
}

Example categories:

Haircut  
ChemicalTreatment  
Coloring  
AddOn  
Retail  
---

# **5\. PRODUCT**

Product can represent both **services and retail items**.

type Product struct {  
	BaseModel

	Name string

	ProductType string  
	// SERVICE | RETAIL

	CategoryID uuid.UUID  
	Category   ProductCategory

	BasePrice int64  
	CostPrice int64

	IsActive bool  
}  
---

# **6\. BRANCH PRODUCT (PRICE \+ STOCK)**

Branch-level price override and inventory.

type BranchProduct struct {  
	BaseModel

	BranchID  uuid.UUID  
	ProductID uuid.UUID

	PriceOverride \*int64

	Stock int64

	Branch  Branch  
	Product Product  
}

Unique constraint:

gorm:"uniqueIndex:idx\_branch\_product"  
---

# **7\. STYLIST**

type Stylist struct {  
	BaseModel

	Name string

	IsActive bool  
}  
---

# **8\. BRANCH STYLIST**

Haircut price override per stylist per branch.

type BranchStylist struct {  
	BaseModel

	BranchID  uuid.UUID  
	StylistID uuid.UUID

	HaircutPriceOverride \*int64

	Branch  Branch  
	Stylist Stylist  
}  
---

# **9\. CUSTOMER**

type Customer struct {  
	BaseModel

	Name  string  
	Phone string

	LoyaltyUserID string  
}  
---

# **10\. AFFILIATE**

type Affiliate struct {  
	BaseModel

	LoyaltyUserID string \`gorm:"uniqueIndex"\`

	AffiliateCode string \`gorm:"uniqueIndex"\`

	Name string

	CommissionType string  
	// PERCENTAGE | FIXED

	CommissionPercentage float64  
	CommissionFixed      int64

	IsActive bool  
}  
---

# **11\. TRANSACTION**

Core POS transaction.

type Transaction struct {  
	BaseModel

	InvoiceNo string \`gorm:"uniqueIndex"\`

	BranchID uuid.UUID  
	Branch   Branch

	CustomerID \*uuid.UUID  
	Customer   \*Customer

	AffiliateID \*uuid.UUID  
	Affiliate   \*Affiliate

	SubtotalAmount int64  
	DiscountAmount int64  
	TotalAmount    int64

	AffiliateCommissionAmountSnapshot int64

	Status string  
	// DRAFT | COMPLETED | VOIDED

	CashDrawerID uuid.UUID  
	CashDrawer   CashDrawer

	EditedByID \*uuid.UUID  
	EditReason string  
}  
---

# **12\. TRANSACTION ITEM**

Financial truth lives here.

type TransactionItem struct {  
	BaseModel

	TransactionID uuid.UUID  
	Transaction   Transaction

	ProductID uuid.UUID  
	Product   Product

	StylistID \*uuid.UUID  
	Stylist   \*Stylist

	ProductNameSnapshot  string  
	ProductTypeSnapshot  string  
	CategoryNameSnapshot string

	IncomeTypeSnapshot string  
	// HAIRCUT | TREATMENT | PRODUCT

	StylistNameSnapshot string

	PriceSnapshot int64  
	Quantity      int64

	GrossSubtotal int64

	ItemDiscount int64  
	NetSubtotal  int64

	CommissionAmountSnapshot int64

	CostPriceSnapshot int64  
}  
---

# **13\. PAYMENT**

Supports split payment.

type Payment struct {  
	BaseModel

	TransactionID uuid.UUID  
	Transaction   Transaction

	Method string  
	// CASH | QRIS | DEBIT | CREDIT

	Amount int64

	ReferenceNo string  
}  
---

# **14\. CASH DRAWER**

type CashDrawer struct {  
	BaseModel

	BranchID uuid.UUID  
	Branch   Branch

	OpenedAt time.Time  
	ClosedAt \*time.Time

	OpeningAmount int64

	ExpectedCash int64  
	CountedCash  int64  
	Variance     int64

	Status string  
	// OPEN | CLOSING | CLOSED

	ClosingSnapshot datatypes.JSON  
}  
---

# **15\. STOCK MOVEMENT**

Inventory ledger.

type StockMovement struct {  
	BaseModel

	BranchID  uuid.UUID  
	ProductID uuid.UUID

	Change int64

	Type string  
	// SALE | RESTOCK | ADJUSTMENT

	ReferenceID uuid.UUID

	Note string  
}  
---

# **16\. SALARY RECORD**

type SalaryRecord struct {  
	BaseModel

	StylistID uuid.UUID  
	BranchID  uuid.UUID

	Month int  
	Year  int

	TotalSales      int64  
	TotalCommission int64

	Status string  
	// GENERATED | PAID  
}

Unique constraint:

stylist\_id \+ branch\_id \+ month \+ year  
---

# **17\. AFFILIATE COMMISSION**

type AffiliateCommission struct {  
	BaseModel

	AffiliateID uuid.UUID  
	TransactionID uuid.UUID

	CommissionAmount int64

	Status string  
	// PENDING | PAID  
}  
---

# **18\. EXPENSE CATEGORY**

type ExpenseCategory struct {  
	BaseModel

	Name string  
	Code string  
}  
---

# **19\. EXPENSE**

type Expense struct {  
	BaseModel

	BranchID uuid.UUID

	CategoryID uuid.UUID

	Description string

	Amount int64

	ExpenseDate time.Time  
}  
---

# **20\. AUDIT LOG**

type AuditLog struct {  
	BaseModel

	Action string

	EntityType string  
	EntityID   uuid.UUID

	PerformedBy uuid.UUID

	Metadata datatypes.JSON  
}  
---

# **ER DIAGRAM SUMMARY**

Branch  
├─ Users  
├─ Transactions  
├─ CashDrawers  
├─ Expenses  
├─ BranchProducts  
└─ BranchStylists

ProductCategory  
└─ Products

Products  
├─ TransactionItems  
├─ BranchProducts  
└─ StockMovements

Transactions  
├─ TransactionItems  
└─ Payments

Stylists  
└─ BranchStylists

Affiliates  
└─ AffiliateCommissions  
---

# **KEY DESIGN PRINCIPLES**

### **Snapshot Accounting**

Financial reports rely only on snapshot fields.

---

### **Drawer Controlled Operations**

Transactions only allowed when drawer OPEN.

---

### **Inventory Ledger**

Stock never updated without StockMovement.

---

### **Atomic Checkout**

Checkout must run inside DB transaction.

