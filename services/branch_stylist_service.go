package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type BranchStylistService struct {
	repo repositories.BranchStylistRepository
}

func NewBranchStylistService(repo repositories.BranchStylistRepository) *BranchStylistService {
	return &BranchStylistService{repo: repo}
}

func (s *BranchStylistService) Create(req dto.CreateBranchStylistRequest) (*models.BranchStylist, error) {
	bs := &models.BranchStylist{
		BranchID:             req.BranchID,
		StylistID:            req.StylistID,
		HaircutPriceOverride: req.HaircutPriceOverride,
		CommissionPercentage: req.CommissionPercentage,
	}
	if err := s.repo.Create(bs); err != nil {
		return nil, err
	}
	return bs, nil
}

func (s *BranchStylistService) GetByBranch(branchID uuid.UUID, offset, limit int) ([]models.BranchStylist, int64, error) {
	return s.repo.FindByBranchID(branchID, offset, limit)
}

func (s *BranchStylistService) GetByID(id uuid.UUID) (*models.BranchStylist, error) {
	return s.repo.FindByID(id)
}

func (s *BranchStylistService) Update(id uuid.UUID, req dto.UpdateBranchStylistRequest) (*models.BranchStylist, error) {
	bs, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.HaircutPriceOverride != nil {
		bs.HaircutPriceOverride = req.HaircutPriceOverride
	}
	if req.CommissionPercentage != nil {
		bs.CommissionPercentage = req.CommissionPercentage
	}
	if err := s.repo.Update(bs); err != nil {
		return nil, err
	}
	return bs, nil
}

func (s *BranchStylistService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
