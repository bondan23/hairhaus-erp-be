package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll(offset, limit int) ([]models.Product, int64, error)
	FindByID(id uuid.UUID) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll(offset, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	r.db.Model(&models.Product{}).Count(&total)
	err := r.db.Preload("Category").Offset(offset).Limit(limit).Find(&products).Error
	return products, total, err
}

func (r *productRepository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}
