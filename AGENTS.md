Berikut contoh **`AGENTS.md`** yang biasanya digunakan sebagai **developer guide untuk AI coding agents dan engineers** yang bekerja di repository backend seperti HAIRHAUS ERP.
Dokumen ini merangkum **arsitektur, aturan bisnis, dan aturan coding** agar semua agen (AI maupun developer) memahami sistem.

---

# AGENTS.md

## HAIRHAUS ERP Backend

This document provides guidelines for **AI agents and engineers** contributing to the HAIRHAUS ERP backend system.

The goal of this file is to ensure consistent architecture, safe financial logic, and predictable code generation.

---

# 1. System Overview

HAIRHAUS ERP is a backend system designed to manage **multi-branch barbershop operations**.

The system supports:

* POS transactions
* Cash drawer reconciliation
* Inventory management
* Affiliate commission
* Stylist salary generation
* Expense management
* Financial reporting

The backend is designed with **financial integrity and transactional safety as primary priorities**.

---

# 2. Technology Stack

Backend services must use the following stack:

| Layer              | Technology     |
| ------------------ | -------------- |
| Language           | Go (Golang)    |
| ORM                | GORM           |
| Database           | PostgreSQL     |
| Migration          | golang-migrate |
| Authentication     | JWT            |
| Caching (optional) | Redis          |

---

# 3. Architecture

The system follows a **layered architecture**.

```
HTTP Router
     ↓
Handler
     ↓
Service Layer
     ↓
Repository Layer
     ↓
GORM ORM
     ↓
PostgreSQL
```

Rules:

* **Business logic must only exist in the Service Layer**
* Handlers should only handle request/response
* Repository should contain database access logic only

---

# 4. System Principles

All agents must respect the following system principles:

1. Financial integrity first
2. Snapshot-based accounting
3. Drawer-controlled operations
4. Multi-branch isolation
5. Strong data consistency
6. Audit-ready financial system

Never modify financial data retroactively.

---

# 5. Database Principles

### Primary Database

PostgreSQL is the primary database.

Recommended production configuration:

```
Storage: 20 GB
RAM: 2–4 GB
CPU: 2 vCPU
```

This configuration safely supports:

```
~1000 transactions per day
```

---

### Expected Data Growth

Typical HAIRHAUS branch activity:

```
50 transactions / day / branch
```

Estimated storage usage:

```
≈ 73 MB transactions per year
≈ 0.22 GB total data per year
```

Recommended storage buffer:

```
20 GB PostgreSQL storage
```

---

# 6. Core Database Tables

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
salary
expenses
expense_categories
audit_logs
```

All operational tables must include:

```
branch_id
```

unless the table is global.

---

# 7. Transaction System

The POS transaction system is the **core financial engine**.

Transaction lifecycle:

```
DRAFT → COMPLETED → VOIDED
```

---

### Transaction Rules

Checkout must run inside **one database transaction** including:

* create transaction
* create transaction items
* create payments
* deduct inventory
* create stock movements

If any step fails, the entire transaction must rollback.

---

### Snapshot Accounting

Transaction items must store snapshot values:

```
product_name_snapshot
product_type_snapshot
category_name_snapshot
income_type_snapshot
stylist_name_snapshot
price_snapshot
gross_subtotal
commission_amount_snapshot
cost_price_snapshot
```

Snapshots ensure **historical financial accuracy**.

Never recompute historical commissions.

---

# 8. Payment System

Supported payment methods:

```
CASH
QRIS
DEBIT
CREDIT
```

Rules:

* Split payments are supported
* Payment total must equal transaction total
* Only CASH payments affect drawer balance

---

# 9. Cash Drawer System

Each branch can only have **one open drawer at a time**.

Drawer lifecycle:

```
OPEN → CLOSING → CLOSED
```

Rules:

* Opening requires `opening_amount`
* Closing records `counted_cash`
* Variance must be calculated
* Closed drawers cannot be edited

---

# 10. Inventory System

Inventory changes must always create a stock movement record.

Movement types:

```
SALE
RESTOCK
ADJUSTMENT
```

Direct stock updates are **not allowed**.

Stock deduction during checkout must use:

```
SELECT ... FOR UPDATE
```

to prevent race conditions.

---

# 11. Idempotency

Checkout API must support idempotency keys.

Example:

```
idempotency_key = txn-unique-123
```

If the same key is submitted twice, the system must return the **original transaction result**.

---

# 12. Salary Generation

Stylist salary is calculated based on transaction commission snapshots.

Salary records are generated per:

```
stylist_id + branch_id + month + year
```

This combination must be unique.

Salary generation cannot be recalculated once generated.

---

# 13. Expense System

Expenses represent operational costs (OPEX).

Examples:

* electricity
* rent
* internet
* supplies
* cashier salary

These values are used for **Net Profit calculation**.

---

# 14. Financial Metrics

The system calculates the following financial metrics.

Revenue:

```
Revenue = SUM(subtotal_amount)
```

COGS:

```
COGS =
commission
+ affiliate commission
+ retail cost
+ discount
```

Gross Profit:

```
Gross Profit = Revenue - COGS
```

Net Profit:

```
Net Profit = Gross Profit - OPEX
```

---

# 15. Security

Security mechanisms include:

* JWT authentication
* Role-based authorization
* Input validation
* SQL injection protection
* Audit logging

Audit logs must be created for:

* transaction edits
* payment edits
* drawer closing
* salary generation

---

# 16. Coding Guidelines

Agents must follow these rules when generating code.

### Naming

Use snake_case for database fields.

Example:

```
transaction_items
price_snapshot
gross_subtotal
```

---

### Go Code Style

Use clear service naming.

Example:

```
TransactionService
InventoryService
DrawerService
SalaryService
```

Repository naming:

```
TransactionRepository
ProductRepository
InventoryRepository
```

---

### Error Handling

All services must return explicit errors.

Example:

```
ErrDrawerClosed
ErrPaymentMismatch
ErrInsufficientStock
```

---

# 17. Future Scaling

If HAIRHAUS expands to many branches:

Recommended upgrades:

* Redis caching read-heavy or non-critical data.
* background job queue for reports
* read replicas for analytics

---

# 18. Important Engineering Rule

Never sacrifice **financial correctness** for convenience.

If a design decision risks:

* financial accuracy
* audit integrity
* transaction consistency

the change must be rejected.

---