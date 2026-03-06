package services

import (
	"encoding/json"
	"errors"
	"time"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashDrawerService struct {
	repo        repositories.CashDrawerRepository
	txnRepo     repositories.TransactionRepository
	auditRepo   repositories.AuditLogRepository
}

func NewCashDrawerService(
	repo repositories.CashDrawerRepository,
	txnRepo repositories.TransactionRepository,
	auditRepo repositories.AuditLogRepository,
) *CashDrawerService {
	return &CashDrawerService{repo: repo, txnRepo: txnRepo, auditRepo: auditRepo}
}

func (s *CashDrawerService) Open(req dto.OpenDrawerRequest, userID uuid.UUID) (*models.CashDrawer, error) {
	// Check for existing open drawer
	existing, err := s.repo.FindOpenByBranch(req.BranchID)
	if err == nil && existing != nil {
		return nil, errors.New("an open drawer already exists for this branch")
	}

	drawer := &models.CashDrawer{
		BranchID:      req.BranchID,
		OpenedAt:      time.Now(),
		OpeningAmount: req.OpeningAmount,
		ExpectedCash:  req.OpeningAmount,
		Status:        models.DrawerStatusOpen,
	}

	if err := s.repo.Create(drawer); err != nil {
		return nil, err
	}
	return drawer, nil
}

func (s *CashDrawerService) Close(drawerID uuid.UUID, req dto.CloseDrawerRequest, userID uuid.UUID) (*models.CashDrawer, error) {
	drawer, err := s.repo.FindByID(drawerID)
	if err != nil {
		return nil, errors.New("drawer not found")
	}
	if drawer.Status != models.DrawerStatusOpen {
		return nil, errors.New("drawer is not open")
	}

	// Calculate expected cash from transactions
	txns, err := s.txnRepo.FindByDrawerID(drawerID)
	if err != nil {
		return nil, err
	}

	var totalCashPayments int64
	for _, txn := range txns {
		for _, p := range txn.Payments {
			if p.Method == models.PaymentMethodCash {
				totalCashPayments += p.Amount
			}
		}
	}

	expectedCash := drawer.OpeningAmount + totalCashPayments
	variance := req.CountedCash - expectedCash
	now := time.Now()

	// Build closing snapshot
	snapshot := map[string]interface{}{
		"total_transactions":  len(txns),
		"total_cash_payments": totalCashPayments,
		"opening_amount":      drawer.OpeningAmount,
		"expected_cash":       expectedCash,
		"counted_cash":        req.CountedCash,
		"variance":            variance,
		"closed_at":           now,
	}
	snapshotJSON, _ := json.Marshal(snapshot)

	drawer.Status = models.DrawerStatusClosed
	drawer.ClosedAt = &now
	drawer.ExpectedCash = expectedCash
	drawer.CountedCash = req.CountedCash
	drawer.Variance = variance
	drawer.ClosingSnapshot = snapshotJSON

	if err := s.repo.Update(drawer); err != nil {
		return nil, err
	}

	// Audit log
	metadata, _ := json.Marshal(snapshot)
	_ = s.auditRepo.Create(&models.AuditLog{
		Action:      "DRAWER_CLOSED",
		EntityType:  "CashDrawer",
		EntityID:    drawerID,
		PerformedBy: userID,
		Metadata:    metadata,
	})

	return drawer, nil
}

func (s *CashDrawerService) GetByID(id uuid.UUID) (*models.CashDrawer, error) {
	return s.repo.FindByID(id)
}

func (s *CashDrawerService) GetByBranch(branchID uuid.UUID, offset, limit int) ([]models.CashDrawer, int64, error) {
	return s.repo.FindByBranchID(branchID, offset, limit)
}

func (s *CashDrawerService) GetOpenDrawer(branchID uuid.UUID) (*models.CashDrawer, error) {
	drawer, err := s.repo.FindOpenByBranch(branchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no open drawer found for this branch")
		}
		return nil, err
	}
	return drawer, nil
}
