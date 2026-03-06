package services

import (
	"errors"
	"fmt"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type InventoryService struct {
	bpRepo *repositories.BranchProductRepository
	smRepo *repositories.StockMovementRepository
}

func NewInventoryService(
	bpRepo *repositories.BranchProductRepository,
	smRepo *repositories.StockMovementRepository,
) *InventoryService {
	return &InventoryService{bpRepo: bpRepo, smRepo: smRepo}
}

func (s *InventoryService) Restock(req dto.RestockRequest) (*models.BranchProduct, error) {
	bp, err := s.bpRepo.FindByBranchAndProduct(req.BranchID, req.ProductID)
	if err != nil {
		return nil, errors.New("branch product not found")
	}

	bp.Stock += req.Quantity
	if err := s.bpRepo.Update(bp); err != nil {
		return nil, err
	}

	sm := &models.StockMovement{
		BranchID:  req.BranchID,
		ProductID: req.ProductID,
		Change:    req.Quantity,
		Type:      models.StockMovementTypeRestock,
		Note:      req.Note,
	}
	if err := s.smRepo.Create(sm); err != nil {
		return nil, err
	}

	return bp, nil
}

func (s *InventoryService) AdjustStock(req dto.AdjustStockRequest) (*models.BranchProduct, error) {
	bp, err := s.bpRepo.FindByBranchAndProduct(req.BranchID, req.ProductID)
	if err != nil {
		return nil, errors.New("branch product not found")
	}

	newStock := bp.Stock + req.Change
	if newStock < 0 {
		return nil, fmt.Errorf("adjustment would result in negative stock (current: %d, change: %d)", bp.Stock, req.Change)
	}

	bp.Stock = newStock
	if err := s.bpRepo.Update(bp); err != nil {
		return nil, err
	}

	sm := &models.StockMovement{
		BranchID:  req.BranchID,
		ProductID: req.ProductID,
		Change:    req.Change,
		Type:      models.StockMovementTypeAdjustment,
		Note:      req.Note,
	}
	if err := s.smRepo.Create(sm); err != nil {
		return nil, err
	}

	return bp, nil
}

func (s *InventoryService) GetMovements(branchID, productID uuid.UUID, offset, limit int) ([]models.StockMovement, int64, error) {
	return s.smRepo.FindByBranchAndProduct(branchID, productID, offset, limit)
}
