package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"
	"hairhaus-pos-be/utils"

	"github.com/google/uuid"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(req dto.CreateUserRequest) (*models.User, error) {
	hashedPin, err := utils.HashPassword(req.Pin)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		EmployeeID:  req.EmployeeID,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Pin:         hashedPin,
		Role:        req.Role,
		BranchID:    req.BranchID,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAll(offset, limit int) ([]models.User, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Update(id uuid.UUID, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.EmployeeID != nil {
		user.EmployeeID = *req.EmployeeID
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.Pin != nil {
		hashedPin, err := utils.HashPassword(*req.Pin)
		if err != nil {
			return nil, err
		}
		user.Pin = hashedPin
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.BranchID != nil {
		user.BranchID = *req.BranchID
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
