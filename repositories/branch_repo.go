package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(branch *models.Branch) error {
	return r.db.Create(branch).Error
}

func (r *BranchRepository) FindAll(offset, limit int) ([]models.Branch, int64, error) {
	var branches []models.Branch
	var total int64
	r.db.Model(&models.Branch{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&branches).Error
	return branches, total, err
}

func (r *BranchRepository) FindByID(id uuid.UUID) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.First(&branch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *BranchRepository) FindByCode(code string) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.First(&branch, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *BranchRepository) Update(branch *models.Branch) error {
	return r.db.Save(branch).Error
}

func (r *BranchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Branch{}, "id = ?", id).Error
}
