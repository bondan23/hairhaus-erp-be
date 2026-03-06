package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type BranchProductService struct {
	repo repositories.BranchProductRepository
}

func NewBranchProductService(repo repositories.BranchProductRepository) *BranchProductService {
	return &BranchProductService{repo: repo}
}

func (s *BranchProductService) Create(req dto.CreateBranchProductRequest) (*models.BranchProduct, error) {
	bp := &models.BranchProduct{
		BranchID:      req.BranchID,
		ProductID:     req.ProductID,
		PriceOverride: req.PriceOverride,
		Stock:         req.Stock,
	}
	if err := s.repo.Create(bp); err != nil {
		return nil, err
	}
	return bp, nil
}

func (s *BranchProductService) GetByBranch(branchID uuid.UUID, offset, limit int) ([]models.BranchProduct, int64, error) {
	return s.repo.FindByBranchID(branchID, offset, limit)
}

func (s *BranchProductService) GetByID(id uuid.UUID) (*models.BranchProduct, error) {
	return s.repo.FindByID(id)
}

func (s *BranchProductService) Update(id uuid.UUID, req dto.UpdateBranchProductRequest) (*models.BranchProduct, error) {
	bp, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.PriceOverride != nil {
		bp.PriceOverride = req.PriceOverride
	}
	if err := s.repo.Update(bp); err != nil {
		return nil, err
	}
	return bp, nil
}

func (s *BranchProductService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
