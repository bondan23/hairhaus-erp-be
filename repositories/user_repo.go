package repositories

import (
	"hairhaus-pos-be/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindAll(offset, limit int) ([]models.User, int64, error)
	FindByID(id uuid.UUID) (*models.User, error)
	FindByPhone(phoneNumber string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Preload("Branch").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Branch").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPhone(phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Branch").First(&user, "phone_number = ?", phoneNumber).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

