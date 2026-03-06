package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BranchProductRepository struct {
	db *gorm.DB
}

func NewBranchProductRepository(db *gorm.DB) *BranchProductRepository {
	return &BranchProductRepository{db: db}
}

func (r *BranchProductRepository) Create(bp *models.BranchProduct) error {
	return r.db.Create(bp).Error
}

func (r *BranchProductRepository) FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.BranchProduct, int64, error) {
	var bps []models.BranchProduct
	var total int64
	r.db.Model(&models.BranchProduct{}).Where("branch_id = ?", branchID).Count(&total)
	err := r.db.Preload("Product").Preload("Product.Category").
		Where("branch_id = ?", branchID).
		Offset(offset).Limit(limit).Find(&bps).Error
	return bps, total, err
}

func (r *BranchProductRepository) FindByID(id uuid.UUID) (*models.BranchProduct, error) {
	var bp models.BranchProduct
	err := r.db.Preload("Product").Preload("Product.Category").First(&bp, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

func (r *BranchProductRepository) FindByBranchAndProduct(branchID, productID uuid.UUID) (*models.BranchProduct, error) {
	var bp models.BranchProduct
	err := r.db.Preload("Product").Preload("Product.Category").
		First(&bp, "branch_id = ? AND product_id = ?", branchID, productID).Error
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

// FindByBranchAndProductForUpdate uses SELECT ... FOR UPDATE for inventory locking.
func (r *BranchProductRepository) FindByBranchAndProductForUpdate(tx *gorm.DB, branchID, productID uuid.UUID) (*models.BranchProduct, error) {
	var bp models.BranchProduct
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Product").Preload("Product.Category").
		First(&bp, "branch_id = ? AND product_id = ?", branchID, productID).Error
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

func (r *BranchProductRepository) Update(bp *models.BranchProduct) error {
	return r.db.Save(bp).Error
}

func (r *BranchProductRepository) UpdateWithTx(tx *gorm.DB, bp *models.BranchProduct) error {
	return tx.Save(bp).Error
}

func (r *BranchProductRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.BranchProduct{}, "id = ?", id).Error
}
