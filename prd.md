# **HAIRHAUS ERP**

# **Product Requirements Document (PRD)**

---

# **1\. PRODUCT OVERVIEW**

## **1.1 Vision**

HAIRHAUS ERP is an enterprise system designed to manage multi-branch barbershop operations with a strong focus on:

* Financial integrity  
* Accurate profit sharing  
* Income stream analytics  
* Inventory control  
* Cash drawer reconciliation  
* Scalable system architecture

This system replaces traditional POS solutions and acts as a **Business Operating System for HAIRHAUS**.

---

# **2\. BUSINESS MODEL**

HAIRHAUS operates with three primary income streams.

### **💇 Haircut Income**

Revenue generated from haircut services.

### **🧪 Treatment Income**

Revenue generated from treatment services such as:

* Chemical Treatment  
* Hair Coloring  
* Add-On Services

### **🛍 Product Sales Income**

Revenue generated from retail product sales.

Examples:

* Pomade  
* Hair tonic  
* Hair care products

---

# **3\. USER ROLES**

## **Cashier**

Permissions:

* Create transactions  
* Add payments  
* View inventory  
* Cannot edit completed transactions  
* Cannot close cash drawer

---

## **Manager**

Permissions:

* Edit transactions (only if drawer is OPEN)  
* Edit payment methods (only if drawer is OPEN)  
* Open and close cash drawers  
* Record expenses  
* View reports

---

## **Admin**

Permissions:

* Full system access  
* Manage affiliates  
* Generate salary calculations  
* Configure system settings  
* Access all financial reports

---

# **4\. CORE SYSTEM MODULES**

---

# **4.1 Multi-Branch Management**

The system must support multiple branches.

Rules:

* Product catalog is global  
* Product pricing can be overridden per branch  
* Inventory stock is managed per branch  
* Haircut pricing can be overridden per stylist per branch  
* Cash drawers operate per branch  
* Salary generation is calculated per stylist per branch

All operational tables must include:

branch\_id

---

# **4.2 Transaction Module (Stylist-First POS)**

Transaction workflow:

1. Select Hair Maker (Stylist)  
2. Select service or product  
3. Add item to receipt  
4. Apply discount (optional)  
5. Apply affiliate code (optional)  
6. Add payment  
7. Complete transaction

Items must be grouped by stylist in the receipt.

Example:

Bondan  
\- Haircut  
\- Creambath

Djarot  
\- Haircut  
\- Creambath

---

# **4.3 Transaction Snapshot Accounting**

Financial data must rely on **snapshot fields** to preserve historical accuracy.

TransactionItem snapshot fields include:

* productNameSnapshot  
* productTypeSnapshot  
* categoryNameSnapshot  
* incomeTypeSnapshot  
* stylistNameSnapshot  
* priceSnapshot  
* grossSubtotal  
* commissionAmountSnapshot  
* costPriceSnapshot

Transaction snapshot fields include:

* subtotalAmount  
* transactionDiscount  
* totalAmount  
* affiliateCommissionAmountSnapshot

Snapshots ensure that financial reports remain valid even if product or commission configurations change in the future.

---

# **4.4 Income Classification**

Income type is stored using:

incomeTypeSnapshot

Enum values:

HAIRCUT  
TREATMENT  
PRODUCT

Income calculations:

Haircut Income

SUM(grossSubtotal WHERE incomeType \= HAIRCUT)

Treatment Income

SUM(grossSubtotal WHERE incomeType \= TREATMENT)

Product Sales Income

SUM(grossSubtotal WHERE incomeType \= PRODUCT)

Total Revenue (Gross Sales):

SUM(subtotalAmount)

---

# **4.5 COGS Definition**

Cost of Goods Sold (COGS) consists of:

1. Stylist profit sharing (commission)  
2. Affiliate commission  
3. Retail product cost  
4. Discounts

COGS formula:

COGS \=  
commission  
\+ affiliate commission  
\+ retail cost  
\+ discount

---

# **4.6 Discount Allocation**

Discounts are applied at the transaction level.

For margin analysis per income stream, discounts must be allocated proportionally.

Formula:

incomeRatio \= income / totalRevenue

allocatedDiscount \= transactionDiscount × incomeRatio

This allocation is calculated only in the reporting layer and is not stored in the database.

---

# **4.7 Gross Margin Calculation**

### **Haircut Gross Margin**

(Haircut Income  
 \- Haircut Discount  
 \- Profit Sharing Haircut)  
/ Haircut Income

---

### **Treatment Gross Margin**

(Treatment Income  
 \- Treatment Discount  
 \- Profit Sharing Treatment)  
/ Treatment Income

---

### **Product Gross Margin**

(Product Income  
 \- Product Discount  
 \- Product Cost)  
/ Product Income

Product cost is calculated from:

costPriceSnapshot × quantity

---

# **4.8 Payment System**

Supported payment methods:

* Cash  
* QRIS  
* Debit Card  
* Credit Card

Rules:

* Split payment is supported  
* Payment total must equal transaction total  
* Payment edits allowed only when drawer is OPEN  
* Payment edits must create an audit log

Only cash payments affect drawer cash balance.

---

# **4.9 Cash Drawer Management**

Cash drawer lifecycle:

OPEN → CLOSING → CLOSED

Opening a drawer requires:

openingAmount

Closing a drawer involves:

1. Aggregating payment totals  
2. Calculating expected cash  
3. Inputting counted cash  
4. Calculating variance  
5. Saving a closing snapshot

Closed drawers cannot be edited.

---

# **4.10 Inventory Management**

Inventory rules:

* Stock is stored per branch  
* Only retail products affect inventory stock  
* Every stock change must create a StockMovement record  
* Direct stock updates without movement logs are prohibited

Stock movement types:

* SALE  
* RESTOCK  
* ADJUSTMENT

---

# **4.11 Affiliate System**

Affiliate system integrates with the loyalty member service.

Affiliate data includes:

* loyaltyMemberId  
* affiliateCode  
* commissionType  
* commissionValue

Affiliate commission is stored in the transaction snapshot and counted as part of COGS.

---

# **4.12 Salary Generation**

Salary is calculated based on stylist commission.

Rules:

* Generated per stylist per branch per month  
* Based only on commission snapshots  
* Cannot be recalculated once generated

Unique constraint:

stylist\_id \+ branch\_id \+ month \+ year

---

# **4.13 Expense Management (OPEX)**

Expenses represent operational costs.

Examples:

* Electricity  
* Rent  
* Cashier salary  
* Internet  
* Supplies

Expenses are used in Net Profit calculations.

# **4.14 Customer Identification & Loyalty Integration**

During POS transactions, the cashier must identify the customer using their phone number.
The system will attempt to locate the customer in the following order:

ERP Customer Database

HAIRHAUS Loyalty Member Service

New Customer Registration

This process ensures that customer data is centralized and optionally integrated with the loyalty system.

Customer Lookup Flow

When the cashier inputs a phone number in the Transaction Screen, the system performs the following steps.

Step 1 — Search ERP Customer Table

The system first checks the internal ERP customer table.

Lookup condition:

customers.phone = input_phone_number
If customer is found

The system will:

populate the customer information

attach the customer to the transaction

No additional steps are required.

Step 2 — Check Loyalty Member API

If the phone number is not found in the ERP customer table, the system will query the HAIRHAUS Loyalty Member API.

Endpoint:

POST /open/auth/check

Input:

phone_number
If loyalty member is found

The system will:

retrieve the loyalty member data

automatically create a new ERP customer record

map the loyalty member ID to the ERP customer

Example fields stored in ERP:

customers
----------------------------
name
phone
gender
loyalty_id

Where:

loyalty_id = loyalty_member.user_id

The newly created ERP customer will then be attached to the transaction.

Step 3 — Customer Not Found Anywhere

If the phone number is not found in both systems, the POS will display a New Customer Form.

Form fields:

Name
Gender
Phone Number (auto-filled)

Additionally, the system will request customer consent for loyalty registration.

Loyalty Registration Consent

The customer must choose whether they want their phone number to be registered in the HAIRHAUS Loyalty Program.

Option presented:

Do you want to register this number for HAIRHAUS Loyalty Membership?

Choices:

YES — Register to Loyalty
NO — Only register in ERP

# If Customer Consents to Loyalty Registration

The following process occurs:

* System sends OTP request to loyalty service
* Customer verifies OTP
* Loyalty member account is created
* ERP creates customer record with loyalty_id reference

Flow:

POS → Loyalty API → OTP Verification (Go routine) → Loyalty Account Created
                                                            ↓
                                                    ERP Customer Created

This flow ensure
* customer data consistency
* seamless loyalty integration
* minimal disruption during POS transactions.


ERP data example:

customers
----------------------------
name
phone
gender
loyalty_id

# If Customer Declines Loyalty Registration

The system will:

* create a customer record only in the ERP database
* leave loyalty_id as NULL

Example:

customers
----------------------------
name
phone
gender
loyalty_id = null

# The customer can still complete the transaction without loyalty integration.

Customer Data Structure

ERP customer table example:

customers
------------------------------------------------
id
name
phone
gender
loyalty_id
created_at
updated_at

Where:

loyalty_id = external loyalty member identifier
Key System Rules
Phone Number is Primary Identifier

Customer lookup is always performed using:

phone_number
Loyalty Integration is Optional

Customers are allowed to:
* use ERP without joining loyalty program

Customer Creation Must Be Idempotent
The system must prevent duplicate customers with the same phone number.
Constraint:

UNIQUE(phone)
---

# **5\. FINANCIAL REPORTING**

---

# **5.1 Revenue**

Revenue \= SUM(subtotalAmount)

---

# **5.2 COGS**

COGS \=  
commission  
\+ affiliate commission  
\+ retail cost  
\+ discount

---

# **5.3 Gross Profit**

Gross Profit \= Revenue \- COGS

---

# **5.4 OPEX**

OPEX \= SUM(expense.amount)

---

# **5.5 Net Profit**

Net Profit \= Gross Profit \- OPEX

**6\. TECHNICAL REQUIREMENTS**  
---

# **6.1 Technology Stack**

Backend

* Go (Golang)  
* GORM ORM  
* golang-migrate

Database

* PostgreSQL

Cache

* Redis (optional but recommended)

Authentication

* JWT

---

# **6.2 Backend Architecture**

Layered architecture:

HTTP Router  
   ↓  
Handler  
   ↓  
Service Layer  
   ↓  
Repository Layer  
   ↓  
GORM  
   ↓  
PostgreSQL

Business logic must exist only in the service layer.

---

# **6.3 Database Migration**

Database schema migrations are managed using:

golang-migrate

Migration files example:

0001\_init\_schema.up.sql  
0001\_init\_schema.down.sql

---

# **7\. CONCURRENCY RULES**

Concurrency rules ensure financial data integrity.

---

## **Checkout Atomicity**

Checkout must run inside a single database transaction including:

* create transaction  
* create transaction items  
* create payments  
* deduct stock

---

## **Inventory Locking**

Stock deduction must use row-level locking.

SELECT ... FOR UPDATE

---

## **Drawer Concurrency**

Only one OPEN drawer per branch is allowed.

Unique constraint enforced at the database level.

---

## **Transaction Edit Lock**

Transaction edits use optimistic locking with:

updated\_at

---

## **Salary Generation Lock**

Salary generation is protected with the unique constraint:

stylist\_id \+ branch\_id \+ month \+ year

---

## **POS Idempotency**

Checkout API must support idempotency keys to prevent duplicate transactions from double taps.

---

# **8\. SECURITY REQUIREMENTS**

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

# **9\. SYSTEM PRINCIPLES**

HAIRHAUS ERP is built based on the following principles:

1. Financial integrity first  
2. Snapshot-based accounting  
3. Drawer-controlled operations  
4. Multi-branch isolation  
5. Strong data consistency  
6. Audit-ready financial system

---

# **10\. MVP SCOPE**

Phase 1

* Multi-branch support  
* POS transaction system  
* Payment system  
* Cash drawer management  
* Inventory tracking  
* Expense management  
* Basic financial reporting  
* Audit logging

Phase 2

* Affiliate system  
* Margin analytics  
* Income stream reports

Phase 3

* Advanced KPI dashboards  
* Branch comparison analytics  
* Business intelligence features

