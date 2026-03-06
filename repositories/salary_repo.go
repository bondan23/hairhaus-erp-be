package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalaryRepository interface {
	Create(salary *models.SalaryRecord) error
	FindByBranch(branchID uuid.UUID, offset, limit int) ([]models.SalaryRecord, int64, error)
	FindByID(id uuid.UUID) (*models.SalaryRecord, error)
	Update(salary *models.SalaryRecord) error
}

type salaryRepository struct {
	db *gorm.DB
}

func NewSalaryRepository(db *gorm.DB) SalaryRepository {
	return &salaryRepository{db: db}
}

func (r *salaryRepository) Create(salary *models.SalaryRecord) error {
	return r.db.Create(salary).Error
}

func (r *salaryRepository) FindByBranch(branchID uuid.UUID, offset, limit int) ([]models.SalaryRecord, int64, error) {
	var records []models.SalaryRecord
	var total int64
	r.db.Model(&models.SalaryRecord{}).Where("branch_id = ?", branchID).Count(&total)
	err := r.db.Preload("Stylist").Preload("Branch").
		Where("branch_id = ?", branchID).
		Order("year DESC, month DESC").
		Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

func (r *salaryRepository) FindByID(id uuid.UUID) (*models.SalaryRecord, error) {
	var record models.SalaryRecord
	err := r.db.Preload("Stylist").Preload("Branch").First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *salaryRepository) Update(salary *models.SalaryRecord) error {
	return r.db.Save(salary).Error
}
