package repositories

import (
	"hairhaus-pos-be/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExpenseCategoryRepository interface {
	Create(cat *models.ExpenseCategory) error
	FindAll(offset, limit int) ([]models.ExpenseCategory, int64, error)
	FindByID(id uuid.UUID) (*models.ExpenseCategory, error)
	Update(cat *models.ExpenseCategory) error
	Delete(id uuid.UUID) error
}

type expenseCategoryRepository struct {
	db *gorm.DB
}

func NewExpenseCategoryRepository(db *gorm.DB) ExpenseCategoryRepository {
	return &expenseCategoryRepository{db: db}
}

func (r *expenseCategoryRepository) Create(cat *models.ExpenseCategory) error {
	return r.db.Create(cat).Error
}

func (r *expenseCategoryRepository) FindAll(offset, limit int) ([]models.ExpenseCategory, int64, error) {
	var cats []models.ExpenseCategory
	var total int64
	r.db.Model(&models.ExpenseCategory{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&cats).Error
	return cats, total, err
}

func (r *expenseCategoryRepository) FindByID(id uuid.UUID) (*models.ExpenseCategory, error) {
	var cat models.ExpenseCategory
	err := r.db.First(&cat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *expenseCategoryRepository) Update(cat *models.ExpenseCategory) error {
	return r.db.Save(cat).Error
}

func (r *expenseCategoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ExpenseCategory{}, "id = ?", id).Error
}

// ExpenseRepository
type ExpenseRepository interface {
	Create(expense *models.Expense) error
	FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.Expense, int64, error)
	FindByID(id uuid.UUID) (*models.Expense, error)
	Update(expense *models.Expense) error
	Delete(id uuid.UUID) error
	SumByBranchAndDateRange(branchID uuid.UUID, start, end time.Time) (int64, error)
}

type expenseRepository struct {
	db *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *expenseRepository) FindByBranchID(branchID uuid.UUID, offset, limit int) ([]models.Expense, int64, error) {
	var expenses []models.Expense
	var total int64
	r.db.Model(&models.Expense{}).Where("branch_id = ?", branchID).Count(&total)
	err := r.db.Preload("Category").Where("branch_id = ?", branchID).
		Order("expense_date DESC").
		Offset(offset).Limit(limit).Find(&expenses).Error
	return expenses, total, err
}

func (r *expenseRepository) FindByID(id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Preload("Category").First(&expense, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *expenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

func (r *expenseRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Expense{}, "id = ?", id).Error
}

// SumByBranchAndDateRange returns total expenses for a branch in a date range.
func (r *expenseRepository) SumByBranchAndDateRange(branchID uuid.UUID, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.Model(&models.Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("branch_id = ? AND expense_date >= ? AND expense_date <= ?", branchID, start, end).
		Scan(&total).Error
	return total, err
}
