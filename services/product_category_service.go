package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type ProductCategoryService struct {
	repo *repositories.ProductCategoryRepository
}

func NewProductCategoryService(repo *repositories.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{repo: repo}
}

func (s *ProductCategoryService) Create(req dto.CreateProductCategoryRequest) (*models.ProductCategory, error) {
	cat := &models.ProductCategory{Name: req.Name, Code: req.Code}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *ProductCategoryService) GetAll(offset, limit int) ([]models.ProductCategory, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *ProductCategoryService) GetByID(id uuid.UUID) (*models.ProductCategory, error) {
	return s.repo.FindByID(id)
}

func (s *ProductCategoryService) Update(id uuid.UUID, req dto.UpdateProductCategoryRequest) (*models.ProductCategory, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Code != nil {
		cat.Code = *req.Code
	}
	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *ProductCategoryService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
