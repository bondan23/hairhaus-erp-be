package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type CustomerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(req dto.CreateCustomerRequest) (*models.Customer, error) {
	customer := &models.Customer{
		Name:              req.Name,
		Phone:             req.Phone,
		Gender:            req.Gender,
		LoyaltyExternalID: req.LoyaltyExternalID,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetAll(offset, limit int) ([]models.Customer, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *CustomerService) GetByID(id uuid.UUID) (*models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomerService) Update(id uuid.UUID, req dto.UpdateCustomerRequest) (*models.Customer, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.LoyaltyExternalID != nil {
		customer.LoyaltyExternalID = *req.LoyaltyExternalID
	}
	if req.Gender != nil {
		customer.Gender = req.Gender
	}
	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
