package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	DB() *gorm.DB
	CreateWithTx(tx *gorm.DB, txn *models.Transaction) error
	CreateItemWithTx(tx *gorm.DB, item *models.TransactionItem) error
	CreatePaymentWithTx(tx *gorm.DB, payment *models.Payment) error
	DeleteItemsWithTx(tx *gorm.DB, txnID uuid.UUID, hardDelete bool) error
	DeletePaymentsWithTx(tx *gorm.DB, txnID uuid.UUID, hardDelete bool) error
	DeleteWithTx(tx *gorm.DB, txnID uuid.UUID, hardDelete bool) error
	FindByID(id uuid.UUID) (*models.Transaction, error)
	FindByIdempotencyKey(key string) (*models.Transaction, error)
	FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.Transaction, int64, error)
	FindByDrawerID(drawerID uuid.UUID) ([]models.Transaction, error)
	Update(txn *models.Transaction) error
	UpdateWithTx(tx *gorm.DB, txn *models.Transaction) error
	CountTodayByBranch(branchID uuid.UUID) (int64, error)
	GetStylistCommissionSummary(stylistID, branchID uuid.UUID, month, year int) (totalSales int64, totalCommission int64, err error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) DB() *gorm.DB {
	return r.db
}

func (r *transactionRepository) CreateWithTx(tx *gorm.DB, txn *models.Transaction) error {
	return tx.Create(txn).Error
}

func (r *transactionRepository) CreateItemWithTx(tx *gorm.DB, item *models.TransactionItem) error {
	return tx.Create(item).Error
}

func (r *transactionRepository) CreatePaymentWithTx(tx *gorm.DB, payment *models.Payment) error {
	return tx.Create(payment).Error
}

func (r *transactionRepository) DeleteItemsWithTx(tx *gorm.DB, txnID uuid.UUID, hardDelete bool) error {
	query := tx.Where("transaction_id = ?", txnID)
	if hardDelete {
		query = query.Unscoped()
	}
	return query.Delete(&models.TransactionItem{}).Error
}

func (r *transactionRepository) DeletePaymentsWithTx(tx *gorm.DB, txnID uuid.UUID, hardDelete bool) error {
	query := tx.Where("transaction_id = ?", txnID)
	if hardDelete {
		query = query.Unscoped()
	}
	return query.Delete(&models.Payment{}).Error
}

func (r *transactionRepository) DeleteWithTx(tx *gorm.DB, txnID uuid.UUID, hardDelete bool) error {
	query := tx.Where("id = ?", txnID)
	if hardDelete {
		query = query.Unscoped()
	}
	return query.Delete(&models.Transaction{}).Error
}

func (r *transactionRepository) FindByID(id uuid.UUID) (*models.Transaction, error) {
	var txn models.Transaction
	err := r.db.Preload("Items").Preload("Payments").Preload("Branch").
		Preload("Customer").Preload("Affiliate").
		First(&txn, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *transactionRepository) FindByIdempotencyKey(key string) (*models.Transaction, error) {
	var txn models.Transaction
	err := r.db.Preload("Items").Preload("Payments").
		First(&txn, "idempotency_key = ?", key).Error
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

func (r *transactionRepository) FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.Transaction, int64, error) {
	var txns []models.Transaction
	var total int64
	r.db.Model(&models.Transaction{}).Where("branch_id = ?", branchID).Count(&total)
	err := r.db.Preload("Items").Preload("Payments").
		Where("branch_id = ?", branchID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&txns).Error
	return txns, total, err
}

func (r *transactionRepository) FindByDrawerID(drawerID uuid.UUID) ([]models.Transaction, error) {
	var txns []models.Transaction
	err := r.db.Preload("Items").Preload("Payments").
		Where("cash_drawer_id = ? AND status = ?", drawerID, models.TransactionStatusCompleted).
		Find(&txns).Error
	return txns, err
}

func (r *transactionRepository) Update(txn *models.Transaction) error {
	return r.db.Save(txn).Error
}

func (r *transactionRepository) UpdateWithTx(tx *gorm.DB, txn *models.Transaction) error {
	return tx.Save(txn).Error
}

// CountTodayByBranch counts today's transactions for a branch (for invoice numbering).
func (r *transactionRepository) CountTodayByBranch(branchID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Transaction{}).
		Where("branch_id = ? AND DATE(created_at) = CURRENT_DATE", branchID).
		Count(&count).Error
	return count, err
}

// GetStylistCommissionSummary gets total sales and commission for a stylist in a branch for a given month/year.
func (r *transactionRepository) GetStylistCommissionSummary(stylistID, branchID uuid.UUID, month, year int) (totalSales int64, totalCommission int64, err error) {
	type Result struct {
		TotalSales      int64
		TotalCommission int64
	}
	var result Result
	err = r.db.Model(&models.TransactionItem{}).
		Select("COALESCE(SUM(gross_subtotal), 0) as total_sales, COALESCE(SUM(commission_amount_snapshot), 0) as total_commission").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transaction_items.stylist_id = ? AND transactions.branch_id = ? AND transactions.status = ? AND EXTRACT(MONTH FROM transactions.created_at) = ? AND EXTRACT(YEAR FROM transactions.created_at) = ?",
			stylistID, branchID, models.TransactionStatusCompleted, month, year).
		Scan(&result).Error
	return result.TotalSales, result.TotalCommission, err
}
