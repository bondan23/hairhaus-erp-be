package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StylistRepository struct {
	db *gorm.DB
}

func NewStylistRepository(db *gorm.DB) *StylistRepository {
	return &StylistRepository{db: db}
}

func (r *StylistRepository) Create(stylist *models.Stylist) error {
	return r.db.Create(stylist).Error
}

func (r *StylistRepository) FindAll(offset, limit int) ([]models.Stylist, int64, error) {
	var stylists []models.Stylist
	var total int64
	r.db.Model(&models.Stylist{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&stylists).Error
	return stylists, total, err
}

func (r *StylistRepository) FindByID(id uuid.UUID) (*models.Stylist, error) {
	var stylist models.Stylist
	err := r.db.First(&stylist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &stylist, nil
}

func (r *StylistRepository) Update(stylist *models.Stylist) error {
	return r.db.Save(stylist).Error
}

func (r *StylistRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Stylist{}, "id = ?", id).Error
}
