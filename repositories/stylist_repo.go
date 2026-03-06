package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StylistRepository interface {
	Create(stylist *models.Stylist) error
	FindAll(offset, limit int) ([]models.Stylist, int64, error)
	FindByID(id uuid.UUID) (*models.Stylist, error)
	Update(stylist *models.Stylist) error
	Delete(id uuid.UUID) error
}

type stylistRepository struct {
	db *gorm.DB
}

func NewStylistRepository(db *gorm.DB) StylistRepository {
	return &stylistRepository{db: db}
}

func (r *stylistRepository) Create(stylist *models.Stylist) error {
	return r.db.Create(stylist).Error
}

func (r *stylistRepository) FindAll(offset, limit int) ([]models.Stylist, int64, error) {
	var stylists []models.Stylist
	var total int64
	r.db.Model(&models.Stylist{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&stylists).Error
	return stylists, total, err
}

func (r *stylistRepository) FindByID(id uuid.UUID) (*models.Stylist, error) {
	var stylist models.Stylist
	err := r.db.First(&stylist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &stylist, nil
}

func (r *stylistRepository) Update(stylist *models.Stylist) error {
	return r.db.Save(stylist).Error
}

func (r *stylistRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Stylist{}, "id = ?", id).Error
}
