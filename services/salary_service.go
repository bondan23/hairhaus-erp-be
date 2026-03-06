package services

import (
	"encoding/json"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type SalaryService struct {
	salaryRepo repositories.SalaryRepository
	txnRepo    repositories.TransactionRepository
	auditRepo  repositories.AuditLogRepository
}

func NewSalaryService(
	salaryRepo repositories.SalaryRepository,
	txnRepo repositories.TransactionRepository,
	auditRepo repositories.AuditLogRepository,
) *SalaryService {
	return &SalaryService{salaryRepo: salaryRepo, txnRepo: txnRepo, auditRepo: auditRepo}
}

func (s *SalaryService) Generate(req dto.GenerateSalaryRequest, userID uuid.UUID) (*models.SalaryRecord, error) {
	// Get commission summary from transactions
	totalSales, totalCommission, err := s.txnRepo.GetStylistCommissionSummary(
		req.StylistID, req.BranchID, req.Month, req.Year,
	)
	if err != nil {
		return nil, err
	}

	salary := &models.SalaryRecord{
		StylistID:       req.StylistID,
		BranchID:        req.BranchID,
		Month:           req.Month,
		Year:            req.Year,
		TotalSales:      totalSales,
		TotalCommission: totalCommission,
		Status:          models.SalaryStatusGenerated,
	}

	if err := s.salaryRepo.Create(salary); err != nil {
		return nil, err // Unique constraint will prevent duplicates
	}

	// Audit log
	metadata, _ := json.Marshal(map[string]interface{}{
		"stylist_id":       req.StylistID,
		"branch_id":        req.BranchID,
		"month":            req.Month,
		"year":             req.Year,
		"total_sales":      totalSales,
		"total_commission": totalCommission,
	})
	_ = s.auditRepo.Create(&models.AuditLog{
		Action:      "SALARY_GENERATED",
		EntityType:  "SalaryRecord",
		EntityID:    salary.ID,
		PerformedBy: userID,
		Metadata:    metadata,
	})

	return salary, nil
}

func (s *SalaryService) GetByBranch(branchID uuid.UUID, offset, limit int) ([]models.SalaryRecord, int64, error) {
	return s.salaryRepo.FindByBranch(branchID, offset, limit)
}

func (s *SalaryService) GetByID(id uuid.UUID) (*models.SalaryRecord, error) {
	return s.salaryRepo.FindByID(id)
}

func (s *SalaryService) MarkPaid(id uuid.UUID, userID uuid.UUID) (*models.SalaryRecord, error) {
	salary, err := s.salaryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	salary.Status = models.SalaryStatusPaid
	if err := s.salaryRepo.Update(salary); err != nil {
		return nil, err
	}
	return salary, nil
}
