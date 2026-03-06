package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashDrawerRepository struct {
	db *gorm.DB
}

func NewCashDrawerRepository(db *gorm.DB) *CashDrawerRepository {
	return &CashDrawerRepository{db: db}
}

func (r *CashDrawerRepository) Create(drawer *models.CashDrawer) error {
	return r.db.Create(drawer).Error
}

func (r *CashDrawerRepository) FindOpenByBranch(branchID uuid.UUID) (*models.CashDrawer, error) {
	var drawer models.CashDrawer
	err := r.db.First(&drawer, "branch_id = ? AND status = ?", branchID, models.DrawerStatusOpen).Error
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

func (r *CashDrawerRepository) FindByID(id uuid.UUID) (*models.CashDrawer, error) {
	var drawer models.CashDrawer
	err := r.db.Preload("Branch").First(&drawer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &drawer, nil
}

func (r *CashDrawerRepository) FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.CashDrawer, int64, error) {
	var drawers []models.CashDrawer
	var total int64
	r.db.Model(&models.CashDrawer{}).Where("branch_id = ?", branchID).Count(&total)
	err := r.db.Where("branch_id = ?", branchID).
		Order("opened_at DESC").
		Offset(offset).Limit(limit).Find(&drawers).Error
	return drawers, total, err
}

func (r *CashDrawerRepository) Update(drawer *models.CashDrawer) error {
	return r.db.Save(drawer).Error
}
