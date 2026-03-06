package services

import (
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(req dto.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		ProductType: req.ProductType,
		CategoryID:  req.CategoryID,
		BasePrice:   req.BasePrice,
		CostPrice:   req.CostPrice,
		IsActive:    true,
	}
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetAll(offset, limit int) ([]models.Product, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *ProductService) GetByID(id uuid.UUID) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Update(id uuid.UUID, req dto.UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.ProductType != nil {
		product.ProductType = *req.ProductType
	}
	if req.CategoryID != nil {
		product.CategoryID = *req.CategoryID
	}
	if req.BasePrice != nil {
		product.BasePrice = *req.BasePrice
	}
	if req.CostPrice != nil {
		product.CostPrice = *req.CostPrice
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
