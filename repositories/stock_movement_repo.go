package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovementRepository interface {
	Create(sm *models.StockMovement) error
	CreateWithTx(tx *gorm.DB, sm *models.StockMovement) error
	FindByBranchAndProduct(branchID, productID uuid.UUID, offset, limit int) ([]models.StockMovement, int64, error)
}

type stockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) StockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) Create(sm *models.StockMovement) error {
	return r.db.Create(sm).Error
}

func (r *stockMovementRepository) CreateWithTx(tx *gorm.DB, sm *models.StockMovement) error {
	return tx.Create(sm).Error
}

func (r *stockMovementRepository) FindByBranchAndProduct(branchID, productID uuid.UUID, offset, limit int) ([]models.StockMovement, int64, error) {
	var sms []models.StockMovement
	var total int64
	r.db.Model(&models.StockMovement{}).
		Where("branch_id = ? AND product_id = ?", branchID, productID).Count(&total)
	err := r.db.Where("branch_id = ? AND product_id = ?", branchID, productID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&sms).Error
	return sms, total, err
}
