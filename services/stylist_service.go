package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type StylistService struct {
	repo repositories.StylistRepository
}

func NewStylistService(repo repositories.StylistRepository) *StylistService {
	return &StylistService{repo: repo}
}

func (s *StylistService) Create(req dto.CreateStylistRequest) (*models.Stylist, error) {
	stylist := &models.Stylist{Name: req.Name, IsActive: true}
	if err := s.repo.Create(stylist); err != nil {
		return nil, err
	}
	return stylist, nil
}

func (s *StylistService) GetAll(offset, limit int) ([]models.Stylist, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *StylistService) GetByID(id uuid.UUID) (*models.Stylist, error) {
	return s.repo.FindByID(id)
}

func (s *StylistService) Update(id uuid.UUID, req dto.UpdateStylistRequest) (*models.Stylist, error) {
	stylist, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		stylist.Name = *req.Name
	}
	if req.IsActive != nil {
		stylist.IsActive = *req.IsActive
	}
	if err := s.repo.Update(stylist); err != nil {
		return nil, err
	}
	return stylist, nil
}

func (s *StylistService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
