package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockMovementRepository struct {
	db *gorm.DB
}

func NewStockMovementRepository(db *gorm.DB) *StockMovementRepository {
	return &StockMovementRepository{db: db}
}

func (r *StockMovementRepository) Create(sm *models.StockMovement) error {
	return r.db.Create(sm).Error
}

func (r *StockMovementRepository) CreateWithTx(tx *gorm.DB, sm *models.StockMovement) error {
	return tx.Create(sm).Error
}

func (r *StockMovementRepository) FindByBranchAndProduct(branchID, productID uuid.UUID, offset, limit int) ([]models.StockMovement, int64, error) {
	var sms []models.StockMovement
	var total int64
	r.db.Model(&models.StockMovement{}).
		Where("branch_id = ? AND product_id = ?", branchID, productID).Count(&total)
	err := r.db.Where("branch_id = ? AND product_id = ?", branchID, productID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&sms).Error
	return sms, total, err
}
