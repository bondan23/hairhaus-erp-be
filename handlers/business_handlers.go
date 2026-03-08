package handlers

import (
	"net/http"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/services"
	"hairhaus-pos-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ===== Cash Drawer Handler =====
type CashDrawerHandler struct {
	service *services.CashDrawerService
}

func NewCashDrawerHandler(service *services.CashDrawerService) *CashDrawerHandler {
	return &CashDrawerHandler{service: service}
}

func (h *CashDrawerHandler) Open(c *gin.Context) {
	var req dto.OpenDrawerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	drawer, err := h.service.Open(req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusConflict, err.Error())
		return
	}
	utils.RespondCreated(c, "Drawer opened", drawer)
}

func (h *CashDrawerHandler) Close(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.CloseDrawerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	drawer, err := h.service.Close(id, req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Drawer closed", drawer)
}

func (h *CashDrawerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	drawer, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Drawer not found")
		return
	}
	utils.RespondSuccess(c, "", drawer)
}

func (h *CashDrawerHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	drawers, total, err := h.service.GetByBranch(branchID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, drawers, page, pageSize, total)
}

func (h *CashDrawerHandler) GetOpen(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	drawer, err := h.service.GetOpenDrawer(branchID)
	if err != nil {
		utils.RespondNotFound(c, err.Error())
		return
	}
	utils.RespondSuccess(c, "", drawer)
}

// ===== Transaction Handler =====
type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Save(c *gin.Context) {
	var req dto.SaveTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	txn, err := h.service.SaveTransaction(req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondCreated(c, "Transaction saved as draft", txn)
}

func (h *TransactionHandler) CheckoutDraft(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid draft ID")
		return
	}

	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	req.TransactionID = &id

	userID := c.MustGet("user_id").(uuid.UUID)
	txn, err := h.service.Checkout(req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Transaction completed", txn)
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	txn, err := h.service.Checkout(req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondCreated(c, "Transaction completed", txn)
}

func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	txn, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Transaction not found")
		return
	}
	utils.RespondSuccess(c, "", txn)
}

func (h *TransactionHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	txns, total, err := h.service.GetByBranch(branchID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, txns, page, pageSize, total)
}

func (h *TransactionHandler) Edit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.EditTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	txn, err := h.service.EditTransaction(id, req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Transaction updated", txn)
}

func (h *TransactionHandler) Void(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var body struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	txn, err := h.service.VoidTransaction(id, userID, body.Reason)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Transaction voided", txn)
}

// ===== Inventory Handler =====
type InventoryHandler struct {
	service *services.InventoryService
}

func NewInventoryHandler(service *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) Restock(c *gin.Context) {
	var req dto.RestockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	bp, err := h.service.Restock(req)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Stock updated", bp)
}

func (h *InventoryHandler) Adjust(c *gin.Context) {
	var req dto.AdjustStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	bp, err := h.service.AdjustStock(req)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Stock adjusted", bp)
}

func (h *InventoryHandler) GetMovements(c *gin.Context) {
	branchID, err := uuid.Parse(c.Query("branch_id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch_id")
		return
	}
	productID, err := uuid.Parse(c.Query("product_id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid product_id")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	movements, total, err := h.service.GetMovements(branchID, productID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, movements, page, pageSize, total)
}

// ===== Affiliate Handler =====
type AffiliateHandler struct {
	service *services.AffiliateService
}

func NewAffiliateHandler(service *services.AffiliateService) *AffiliateHandler {
	return &AffiliateHandler{service: service}
}

func (h *AffiliateHandler) Create(c *gin.Context) {
	var req dto.CreateAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	affiliate, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Affiliate created", affiliate)
}

func (h *AffiliateHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	affiliates, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, affiliates, page, pageSize, total)
}

func (h *AffiliateHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	affiliate, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Affiliate not found")
		return
	}
	utils.RespondSuccess(c, "", affiliate)
}

func (h *AffiliateHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateAffiliateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	affiliate, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Affiliate updated", affiliate)
}

func (h *AffiliateHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Affiliate deleted", nil)
}

func (h *AffiliateHandler) GetCommissions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	comms, total, err := h.service.GetCommissions(id, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, comms, page, pageSize, total)
}

// ===== Salary Handler =====
type SalaryHandler struct {
	service *services.SalaryService
}

func NewSalaryHandler(service *services.SalaryService) *SalaryHandler {
	return &SalaryHandler{service: service}
}

func (h *SalaryHandler) Generate(c *gin.Context) {
	var req dto.GenerateSalaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	salary, err := h.service.Generate(req, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondCreated(c, "Salary generated", salary)
}

func (h *SalaryHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	salaries, total, err := h.service.GetByBranch(branchID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, salaries, page, pageSize, total)
}

func (h *SalaryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	salary, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Salary record not found")
		return
	}
	utils.RespondSuccess(c, "", salary)
}

func (h *SalaryHandler) MarkPaid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	userID := c.MustGet("user_id").(uuid.UUID)
	salary, err := h.service.MarkPaid(id, userID)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "Salary marked as paid", salary)
}

// ===== Expense Category Handler =====
type ExpenseCategoryHandler struct {
	service *services.ExpenseCategoryService
}

func NewExpenseCategoryHandler(service *services.ExpenseCategoryService) *ExpenseCategoryHandler {
	return &ExpenseCategoryHandler{service: service}
}

func (h *ExpenseCategoryHandler) Create(c *gin.Context) {
	var req dto.CreateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	cat, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Expense category created", cat)
}

func (h *ExpenseCategoryHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	cats, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, cats, page, pageSize, total)
}

func (h *ExpenseCategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	cat, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Expense category not found")
		return
	}
	utils.RespondSuccess(c, "", cat)
}

func (h *ExpenseCategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateExpenseCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	cat, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Expense category updated", cat)
}

func (h *ExpenseCategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Expense category deleted", nil)
}

// ===== Expense Handler =====
type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(service *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) Create(c *gin.Context) {
	var req dto.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	expense, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondCreated(c, "Expense created", expense)
}

func (h *ExpenseHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	expenses, total, err := h.service.GetByBranch(branchID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, expenses, page, pageSize, total)
}

func (h *ExpenseHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	expense, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Expense not found")
		return
	}
	utils.RespondSuccess(c, "", expense)
}

func (h *ExpenseHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	expense, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Expense updated", expense)
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Expense deleted", nil)
}

// ===== Report Handler =====
type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) GetFinancialReport(c *gin.Context) {
	var filter dto.ReportFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	report, err := h.service.GetFinancialReport(filter)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "", report)
}

func (h *ReportHandler) GetRevenueByIncomeType(c *gin.Context) {
	var filter dto.ReportFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	breakdown, err := h.service.GetRevenueByIncomeType(filter)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, "", breakdown)
}
