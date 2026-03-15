Berikut **revisi `ARCHITECTURE.md`** dengan **Redis layer ditambahkan sebelum Database Layer**, serta penjelasan bahwa Redis **tidak digunakan untuk semua use case** dan hanya dipakai pada skenario tertentu.

Saya juga menambahkan **section khusus Redis** agar engineer jelas kapan harus memakai cache dan kapan tidak.

---

# ARCHITECTURE.md

## HAIRHAUS ERP Backend Architecture

This document describes the **system architecture of the HAIRHAUS ERP backend**, including system layers, module responsibilities, data flow, and key design principles.

The backend system powers **multi-branch barbershop operations**, including POS transactions, inventory control, financial reporting, and operational analytics.

---

# 1. System Overview

HAIRHAUS ERP is designed as a **transaction-safe financial system** for barbershop operations.

Core capabilities:

* POS transaction processing
* Cash drawer reconciliation
* Inventory tracking
* Affiliate commission tracking
* Stylist salary generation
* Expense management
* Financial reporting

The system prioritizes:

* financial integrity
* consistency
* auditability
* operational simplicity

---

# 2. High Level Architecture

The system follows a **modular layered architecture**.

```
Client Applications
│
├─ POS Frontend
├─ Admin Dashboard
└─ Mobile Applications
        │
        ▼
HTTP API Layer
        │
        ▼
Application Layer (Services)
        │
        ▼
Repository Layer
        │
        ▼
Cache Layer (Redis - Optional)
        │
        ▼
Database Layer (PostgreSQL)
```

Redis is used as an **optional cache layer** and is **not required for all use cases**.

---

# 3. Backend Layered Architecture

Backend services follow a strict layered design.

```
Router
  │
  ▼
Handler (HTTP Layer)
  │
  ▼
Service Layer (Business Logic)
  │
  ▼
Repository Layer
  │
  ▼
Redis Cache (Optional)
  │
  ▼
PostgreSQL Database
```

---

## Handler Layer

Responsibilities:

* parse HTTP request
* validate request input
* call service layer
* format HTTP response

Handlers **must not contain business logic**.

Example endpoint:

```
POST /transactions/checkout
```

Handler calls:

```
TransactionService.Checkout()
```

---

## Service Layer

The service layer contains **all business rules**.

Responsibilities:

* transaction validation
* financial calculations
* commission calculations
* discount allocation
* inventory deduction
* audit logging

All business rules must be implemented here.

Example services:

```
TransactionService
InventoryService
DrawerService
AffiliateService
SalaryService
ExpenseService
ReportService
```

---

## Repository Layer

Repository handles **data access abstraction**.

Responsibilities:

* database queries
* cache lookup logic
* entity persistence

Repositories may interact with:

* Redis (optional)
* PostgreSQL

---

# 4. Redis Cache Layer

Redis is used as an **optional performance optimization layer**.

Redis is **not part of the critical financial data path** and must never be used as the source of truth.

PostgreSQL remains the **primary system of record**.

---

## Redis Use Cases

Redis should be used only for **read-heavy or non-critical data**.

Examples:

### Product Catalog Cache

```
GET /products
GET /branches/{id}/products
```

Product data rarely changes and can be cached.

---

### Branch Configuration

Example cached entities:

```
branch_products
branch_stylists
```

Used frequently by POS systems.

---

### Session / Authentication Cache

Redis can optionally store:

```
JWT session blacklist
temporary auth sessions
```

---

### Rate Limiting

Redis can support:

```
login rate limit
API rate limiting
```

---

### Report Cache

Heavy financial reports may be cached temporarily.

Example:

```
financial reports
income summaries
dashboard analytics
```

Cache TTL example:

```
60 seconds
```

---

# 5. Redis Restrictions

Redis **must NOT be used for the following**:

### POS Transactions

Never cache:

```
transactions
transaction_items
payments
```

These require **strong consistency**.

---

### Inventory Updates

Never cache stock deduction logic.

Inventory must always use **database row locking**.

---

### Cash Drawer State

Drawer status must always come from **PostgreSQL**.

---

# 6. Core System Modules

The HAIRHAUS ERP backend is composed of several functional modules.

---

# 6.1 POS Transaction Engine

The POS engine is the **most critical component of the system**.

Main entities:

```
transactions
transaction_items
payments
cash_drawers
```

Transaction flow:

```
create transaction
    ↓
create transaction items
    ↓
validate payments
    ↓
deduct inventory
    ↓
create stock movements
    ↓
finalize transaction
```

All steps must run inside a **single database transaction**.

---

# 6.2 Inventory Management

Inventory is controlled through **stock movements**.

Key tables:

```
products
branch_products
stock_movements
```

Movement types:

```
SALE
RESTOCK
ADJUSTMENT
```

Stock updates must always create a movement record.

---

# 6.3 Cash Drawer System

Cash drawer ensures **cash reconciliation integrity**.

Drawer lifecycle:

```
OPEN → CLOSING → CLOSED
```

Rules:

* only one open drawer per branch
* transactions require open drawer
* closing records expected vs counted cash

---

# 6.4 Customer Management

Customers are identified primarily by **phone number**.

Customer sources:

```
ERP Database
Loyalty Service
New Registration
```

Phone numbers must be unique.

---

# 6.5 Affiliate System

The affiliate system enables **referral-based commissions**.

Entities:

```
affiliates
affiliate_commissions
```

Affiliate commissions are recorded as part of **COGS**.

---

# 6.6 Salary System

Stylist salary is generated from **commission snapshots stored in transaction items**.

Salary records:

```
salaries
```

Generated per:

```
stylist_id
branch_id
month
year
```

Salary cannot be regenerated once created.

---

# 6.7 Expense Management

Expenses represent **operational costs (OPEX)**.

Entities:

```
expense_categories
expenses
```

Expenses contribute to **net profit calculation**.

---

# 7. Database Architecture

Primary database:

```
PostgreSQL
```

Key characteristics:

* relational data model
* snapshot accounting
* transactional integrity
* strong consistency

---

## Core Tables

Main entities:

```
branches
users
products
product_categories
branch_products
stylists
branch_stylists
customers
cash_drawers
transactions
transaction_items
payments
stock_movements
affiliates
affiliate_commissions
salaries
expenses
expense_categories
audit_logs
```

---

# 8. Transaction Safety

All POS checkout operations must run within a **database transaction**.

Example:

```
BEGIN

create transaction
create transaction items
create payments
update inventory
create stock movements

COMMIT
```

If any step fails:

```
ROLLBACK
```

---

# 9. Concurrency Control

### Inventory Locking

```
SELECT ... FOR UPDATE
```

Prevents overselling inventory.

---

### Drawer Lock

Only one open drawer per branch.

Enforced by database constraint.

---

### Transaction Edit Lock

Uses optimistic locking:

```
updated_at
```

---

### Idempotency

Checkout API must support idempotency keys to prevent duplicate transactions.

---

# 10. Performance Considerations

Typical HAIRHAUS branch activity:

```
50 transactions/day
```

Estimated storage growth:

```
≈ 73MB per year per branch
```

Recommended database storage:

```
20GB PostgreSQL storage
```

---

# 11. Future Scalability

As HAIRHAUS expands, the system can scale with:

* PostgreSQL table partitioning
* Redis caching layer expansion
* read replicas for analytics
* background job workers

---

# 12. Design Philosophy

The HAIRHAUS ERP backend follows these principles:

1. Financial integrity above all
2. Transaction safety
3. Audit-ready system design
4. Clear separation of concerns
5. Predictable system behavior

Every architectural decision must respect these principles.

---