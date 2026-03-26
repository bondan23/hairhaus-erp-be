package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	FindAll(offset, limit int, name, phone string) ([]models.Customer, int64, error)
	FindByID(id uuid.UUID) (*models.Customer, error)
	FindByPhone(phone string) (*models.Customer, error)
	Update(customer *models.Customer) error
	Delete(id uuid.UUID) error
	HardDelete(id uuid.UUID) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) FindAll(offset, limit int, name, phone string) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64
	query := r.db.Model(&models.Customer{})
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Find(&customers).Error
	return customers, total, err
}

func (r *customerRepository) FindByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByPhone(phone string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "phone = ?", phone).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) HardDelete(id uuid.UUID) error {
	return r.db.Unscoped().Delete(&models.Customer{}, "id = ?", id).Error
}
