package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategoryRepository interface {
	Create(cat *models.ProductCategory) error
	FindAll(offset, limit int) ([]models.ProductCategory, int64, error)
	FindByID(id uuid.UUID) (*models.ProductCategory, error)
	Update(cat *models.ProductCategory) error
	Delete(id uuid.UUID) error
}

type productCategoryRepository struct {
	db *gorm.DB
}

func NewProductCategoryRepository(db *gorm.DB) ProductCategoryRepository {
	return &productCategoryRepository{db: db}
}

func (r *productCategoryRepository) Create(cat *models.ProductCategory) error {
	return r.db.Create(cat).Error
}

func (r *productCategoryRepository) FindAll(offset, limit int) ([]models.ProductCategory, int64, error) {
	var cats []models.ProductCategory
	var total int64
	r.db.Model(&models.ProductCategory{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&cats).Error
	return cats, total, err
}

func (r *productCategoryRepository) FindByID(id uuid.UUID) (*models.ProductCategory, error) {
	var cat models.ProductCategory
	err := r.db.First(&cat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *productCategoryRepository) Update(cat *models.ProductCategory) error {
	return r.db.Save(cat).Error
}

func (r *productCategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ProductCategory{}, "id = ?", id).Error
}
