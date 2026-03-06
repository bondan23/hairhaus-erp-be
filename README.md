# HAIRHAUS ERP Backend — Walkthrough

## Overview

Built a complete Go backend service for HAIRHAUS ERP, a multi-branch barbershop management system. The service covers all MVP Phase 1 + Phase 2 requirements from the PRD.

## Project Structure

```
hairhaus-pos-be/
├── main.go                  # Entry point, DI wiring
├── config/
│   ├── config.go            # Env-based config
│   ├── database.go          # GORM + PostgreSQL
│   └── redis.go             # Redis client
├── models/                  # 20 GORM models
├── repositories/            # Data access layer
├── services/                # Business logic layer
├── handlers/                # HTTP handlers
├── routes/routes.go         # Route registration + middleware
├── middleware/auth.go       # JWT auth + role-based access
├── dto/dto.go               # Request/response DTOs
├── utils/                   # JWT, password, response, pagination, invoice
└── migrations/              # SQL migration files
```

## Key Design Decisions

| Decision | Rationale |
|---|---|
| **Snapshot accounting** | All financial fields captured at checkout time, preserving historical accuracy |
| **Atomic checkout** | Single DB transaction for creating transaction, items, payments, and stock deduction |
| **`SELECT ... FOR UPDATE`** | Row-level locking on inventory prevents overselling |
| **Idempotency keys** | Prevents duplicate transactions from double-taps |
| **Role-based middleware** | ADMIN, MANAGER, CASHIER permissions enforced per route |
| **One open drawer per branch** | Checked at application level when opening a drawer |
| **Commission from snapshots** | Salary generation reads `commission_amount_snapshot` from completed transactions |

## API Endpoints

### Public
| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/auth/login` | JWT login |

### Protected (requires Bearer token)

**CRUD Modules**: Branches, Users, Product Categories, Products, Stylists, Branch Products, Branch Stylists, Customers

**Business Modules**:
| Method | Endpoint | Roles | Description |
|---|---|---|---|
| POST | `/api/v1/drawers/open` | ADMIN, MANAGER | Open cash drawer |
| POST | `/api/v1/drawers/:id/close` | ADMIN, MANAGER | Close with variance calc |
| POST | `/api/v1/transactions/checkout` | ALL | Atomic POS checkout |
| PUT | `/api/v1/transactions/:id` | ADMIN, MANAGER | Edit (drawer must be open) |
| POST | `/api/v1/transactions/:id/void` | ADMIN, MANAGER | Void transaction |
| POST | `/api/v1/inventory/restock` | ADMIN, MANAGER | Restock inventory |
| POST | `/api/v1/inventory/adjust` | ADMIN, MANAGER | Stock adjustment |
| POST | `/api/v1/salaries/generate` | ADMIN | Generate monthly salary |
| GET | `/api/v1/reports/financial` | ADMIN, MANAGER | Revenue, COGS, GP, OPEX, NP |
| GET | `/api/v1/reports/income-by-type` | ADMIN, MANAGER | Haircut/Treatment/Product breakdown |

## Verification

| Check | Result |
|---|---|
| `go build ./...` | ✅ Pass |
| `go vet ./...` | ✅ Pass |

## Getting Started

1. Copy `.env.example` to `.env` and configure PostgreSQL/Redis credentials
2. Run `go run main.go` — the server auto-migrates all tables and starts on the configured port
3. Create an admin user (first time via direct DB insert or seed script) to access protected endpoints
