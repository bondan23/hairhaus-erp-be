package services

import (
	"errors"
	"time"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/google/uuid"
)

type ExpenseCategoryService struct {
	repo *repositories.ExpenseCategoryRepository
}

func NewExpenseCategoryService(repo *repositories.ExpenseCategoryRepository) *ExpenseCategoryService {
	return &ExpenseCategoryService{repo: repo}
}

func (s *ExpenseCategoryService) Create(req dto.CreateExpenseCategoryRequest) (*models.ExpenseCategory, error) {
	cat := &models.ExpenseCategory{Name: req.Name, Code: req.Code}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *ExpenseCategoryService) GetAll(offset, limit int) ([]models.ExpenseCategory, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *ExpenseCategoryService) GetByID(id uuid.UUID) (*models.ExpenseCategory, error) {
	return s.repo.FindByID(id)
}

func (s *ExpenseCategoryService) Update(id uuid.UUID, req dto.UpdateExpenseCategoryRequest) (*models.ExpenseCategory, error) {
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

func (s *ExpenseCategoryService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// ExpenseService
type ExpenseService struct {
	repo *repositories.ExpenseRepository
}

func NewExpenseService(repo *repositories.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

func (s *ExpenseService) Create(req dto.CreateExpenseRequest) (*models.Expense, error) {
	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		return nil, errors.New("invalid expense_date format, use YYYY-MM-DD")
	}

	expense := &models.Expense{
		BranchID:    req.BranchID,
		CategoryID:  req.CategoryID,
		Description: req.Description,
		Amount:      req.Amount,
		ExpenseDate: expenseDate,
	}
	if err := s.repo.Create(expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) GetByBranch(branchID uuid.UUID, offset, limit int) ([]models.Expense, int64, error) {
	return s.repo.FindByBranchID(branchID, offset, limit)
}

func (s *ExpenseService) GetByID(id uuid.UUID) (*models.Expense, error) {
	return s.repo.FindByID(id)
}

func (s *ExpenseService) Update(id uuid.UUID, req dto.UpdateExpenseRequest) (*models.Expense, error) {
	expense, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.CategoryID != nil {
		expense.CategoryID = *req.CategoryID
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.ExpenseDate != nil {
		date, err := time.Parse("2006-01-02", *req.ExpenseDate)
		if err != nil {
			return nil, errors.New("invalid expense_date format")
		}
		expense.ExpenseDate = date
	}
	if err := s.repo.Update(expense); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *ExpenseService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
