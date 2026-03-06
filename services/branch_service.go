package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type BranchService struct {
	repo *repositories.BranchRepository
}

func NewBranchService(repo *repositories.BranchRepository) *BranchService {
	return &BranchService{repo: repo}
}

func (s *BranchService) Create(req dto.CreateBranchRequest) (*models.Branch, error) {
	branch := &models.Branch{
		Name:     req.Name,
		Code:     req.Code,
		Address:  req.Address,
		Phone:    req.Phone,
		IsActive: true,
	}
	if err := s.repo.Create(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *BranchService) GetAll(offset, limit int) ([]models.Branch, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *BranchService) GetByID(id uuid.UUID) (*models.Branch, error) {
	return s.repo.FindByID(id)
}

func (s *BranchService) Update(id uuid.UUID, req dto.UpdateBranchRequest) (*models.Branch, error) {
	branch, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Address != nil {
		branch.Address = *req.Address
	}
	if req.Phone != nil {
		branch.Phone = *req.Phone
	}
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}
	if err := s.repo.Update(branch); err != nil {
		return nil, err
	}
	return branch, nil
}

func (s *BranchService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
