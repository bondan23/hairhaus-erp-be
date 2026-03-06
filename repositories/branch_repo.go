package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BranchRepository interface {
	Create(branch *models.Branch) error
	FindAll(offset, limit int) ([]models.Branch, int64, error)
	FindByID(id uuid.UUID) (*models.Branch, error)
	FindByCode(code string) (*models.Branch, error)
	Update(branch *models.Branch) error
	Delete(id uuid.UUID) error
}

type branchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(branch *models.Branch) error {
	return r.db.Create(branch).Error
}

func (r *branchRepository) FindAll(offset, limit int) ([]models.Branch, int64, error) {
	var branches []models.Branch
	var total int64
	r.db.Model(&models.Branch{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&branches).Error
	return branches, total, err
}

func (r *branchRepository) FindByID(id uuid.UUID) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.First(&branch, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) FindByCode(code string) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.First(&branch, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) Update(branch *models.Branch) error {
	return r.db.Save(branch).Error
}

func (r *branchRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Branch{}, "id = ?", id).Error
}
