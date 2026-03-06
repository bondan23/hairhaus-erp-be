package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type AffiliateService struct {
	repo     repositories.AffiliateRepository
	commRepo repositories.AffiliateCommissionRepository
}

func NewAffiliateService(repo repositories.AffiliateRepository, commRepo repositories.AffiliateCommissionRepository) *AffiliateService {
	return &AffiliateService{repo: repo, commRepo: commRepo}
}

func (s *AffiliateService) Create(req dto.CreateAffiliateRequest) (*models.Affiliate, error) {
	affiliate := &models.Affiliate{
		LoyaltyMemberID:      req.LoyaltyMemberID,
		AffiliateCode:        req.AffiliateCode,
		Name:                 req.Name,
		CommissionType:       req.CommissionType,
		CommissionPercentage: req.CommissionPercentage,
		CommissionFixed:      req.CommissionFixed,
		IsActive:             true,
	}
	if err := s.repo.Create(affiliate); err != nil {
		return nil, err
	}
	return affiliate, nil
}

func (s *AffiliateService) GetAll(offset, limit int) ([]models.Affiliate, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *AffiliateService) GetByID(id uuid.UUID) (*models.Affiliate, error) {
	return s.repo.FindByID(id)
}

func (s *AffiliateService) Update(id uuid.UUID, req dto.UpdateAffiliateRequest) (*models.Affiliate, error) {
	affiliate, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		affiliate.Name = *req.Name
	}
	if req.CommissionType != nil {
		affiliate.CommissionType = *req.CommissionType
	}
	if req.CommissionPercentage != nil {
		affiliate.CommissionPercentage = *req.CommissionPercentage
	}
	if req.CommissionFixed != nil {
		affiliate.CommissionFixed = *req.CommissionFixed
	}
	if req.IsActive != nil {
		affiliate.IsActive = *req.IsActive
	}
	if err := s.repo.Update(affiliate); err != nil {
		return nil, err
	}
	return affiliate, nil
}

func (s *AffiliateService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *AffiliateService) GetCommissions(affiliateID uuid.UUID, offset, limit int) ([]models.AffiliateCommission, int64, error) {
	return s.commRepo.FindByAffiliateID(affiliateID, offset, limit)
}
