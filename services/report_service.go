package services

import (
	"errors"
	"time"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportService struct {
	txnRepo     repositories.TransactionRepository
	expenseRepo repositories.ExpenseRepository
}

func NewReportService(txnRepo repositories.TransactionRepository, expenseRepo repositories.ExpenseRepository) *ReportService {
	return &ReportService{txnRepo: txnRepo, expenseRepo: expenseRepo}
}

func (s *ReportService) GetFinancialReport(filter dto.ReportFilter) (*dto.FinancialReport, error) {
	if filter.BranchID == uuid.Nil {
		return nil, errors.New("branch_id is required")
	}

	startDate, err := time.Parse("2006-01-02", filter.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format")
	}
	endDate, err := time.Parse("2006-01-02", filter.EndDate)
	if err != nil {
		return nil, errors.New("invalid end_date format")
	}
	endDate = endDate.Add(24*time.Hour - time.Second) // End of day

	db := s.txnRepo.DB()

	// Revenue = SUM(subtotal_amount) for completed transactions
	var revenue int64
	db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(subtotal_amount), 0)").
		Where("branch_id = ? AND status = ? AND created_at >= ? AND created_at <= ?",
			filter.BranchID, models.TransactionStatusCompleted, startDate, endDate).
		Scan(&revenue)

	// COGS components from transaction items
	type COGSResult struct {
		TotalCommission int64
		TotalRetailCost int64
	}
	var cogsResult COGSResult
	db.Model(&models.TransactionItem{}).
		Select(`COALESCE(SUM(commission_amount_snapshot), 0) as total_commission,
		        COALESCE(SUM(CASE WHEN product_type_snapshot = 'RETAIL' THEN cost_price_snapshot * quantity ELSE 0 END), 0) as total_retail_cost`).
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.branch_id = ? AND transactions.status = ? AND transactions.created_at >= ? AND transactions.created_at <= ?",
			filter.BranchID, models.TransactionStatusCompleted, startDate, endDate).
		Scan(&cogsResult)

	// Affiliate commission
	var affiliateCommission int64
	db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(affiliate_commission_amount_snapshot), 0)").
		Where("branch_id = ? AND status = ? AND created_at >= ? AND created_at <= ?",
			filter.BranchID, models.TransactionStatusCompleted, startDate, endDate).
		Scan(&affiliateCommission)

	// Discount
	var totalDiscount int64
	db.Model(&models.Transaction{}).
		Select("COALESCE(SUM(discount_amount), 0)").
		Where("branch_id = ? AND status = ? AND created_at >= ? AND created_at <= ?",
			filter.BranchID, models.TransactionStatusCompleted, startDate, endDate).
		Scan(&totalDiscount)

	cogsTotal := cogsResult.TotalCommission + affiliateCommission + cogsResult.TotalRetailCost + totalDiscount

	// OPEX
	opex, _ := s.expenseRepo.SumByBranchAndDateRange(filter.BranchID, startDate, endDate)

	grossProfit := revenue - cogsTotal
	netProfit := grossProfit - opex

	return &dto.FinancialReport{
		Revenue: revenue,
		COGS: dto.COGSBreakdown{
			Commission:          cogsResult.TotalCommission,
			AffiliateCommission: affiliateCommission,
			RetailCost:          cogsResult.TotalRetailCost,
			Discount:            totalDiscount,
			Total:               cogsTotal,
		},
		GrossProfit: grossProfit,
		OPEX:        opex,
		NetProfit:   netProfit,
	}, nil
}

// GetRevenueByIncomeType returns revenue broken down by income type.
func (s *ReportService) GetRevenueByIncomeType(filter dto.ReportFilter) (map[string]int64, error) {
	if filter.BranchID == uuid.Nil {
		return nil, errors.New("branch_id is required")
	}

	startDate, err := time.Parse("2006-01-02", filter.StartDate)
	if err != nil {
		return nil, errors.New("invalid start_date format")
	}
	endDate, err := time.Parse("2006-01-02", filter.EndDate)
	if err != nil {
		return nil, errors.New("invalid end_date format")
	}
	endDate = endDate.Add(24*time.Hour - time.Second)

	db := s.txnRepo.DB()

	type IncomeRow struct {
		IncomeTypeSnapshot string
		Total              int64
	}
	var results []IncomeRow
	db.Model(&models.TransactionItem{}).
		Select("income_type_snapshot, COALESCE(SUM(gross_subtotal), 0) as total").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.branch_id = ? AND transactions.status = ? AND transactions.created_at >= ? AND transactions.created_at <= ?",
			filter.BranchID, models.TransactionStatusCompleted, startDate, endDate).
		Group("income_type_snapshot").
		Scan(&results)

	breakdown := map[string]int64{
		models.IncomeTypeHaircut:   0,
		models.IncomeTypeTreatment: 0,
		models.IncomeTypeProduct:   0,
	}
	for _, r := range results {
		breakdown[r.IncomeTypeSnapshot] = r.Total
	}

	return breakdown, nil
}

// Ensure gorm.DB is importable
var _ *gorm.DB
