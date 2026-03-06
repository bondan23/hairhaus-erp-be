package handlers

import (
	"net/http"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/services"
	"hairhaus-pos-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	resp, err := h.service.Login(req)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.RespondSuccess(c, "Login successful", resp)
}

// ===== Branch Handler =====
type BranchHandler struct {
	service *services.BranchService
}

func NewBranchHandler(service *services.BranchService) *BranchHandler {
	return &BranchHandler{service: service}
}

func (h *BranchHandler) Create(c *gin.Context) {
	var req dto.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	branch, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Branch created", branch)
}

func (h *BranchHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	branches, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, branches, page, pageSize, total)
}

func (h *BranchHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	branch, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Branch not found")
		return
	}
	utils.RespondSuccess(c, "", branch)
}

func (h *BranchHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	branch, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Branch updated", branch)
}

func (h *BranchHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Branch deleted", nil)
}

// ===== User Handler =====
type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	user, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "User created", user)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	users, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, users, page, pageSize, total)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	user, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "User not found")
		return
	}
	utils.RespondSuccess(c, "", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	user, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "User updated", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "User deleted", nil)
}

// ===== Product Category Handler =====
type ProductCategoryHandler struct {
	service *services.ProductCategoryService
}

func NewProductCategoryHandler(service *services.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{service: service}
}

func (h *ProductCategoryHandler) Create(c *gin.Context) {
	var req dto.CreateProductCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	cat, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Product category created", cat)
}

func (h *ProductCategoryHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	cats, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, cats, page, pageSize, total)
}

func (h *ProductCategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	cat, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Product category not found")
		return
	}
	utils.RespondSuccess(c, "", cat)
}

func (h *ProductCategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateProductCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	cat, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Product category updated", cat)
}

func (h *ProductCategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Product category deleted", nil)
}

// ===== Product Handler =====
type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	product, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Product created", product)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	products, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, products, page, pageSize, total)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	product, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Product not found")
		return
	}
	utils.RespondSuccess(c, "", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	product, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Product updated", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Product deleted", nil)
}

// ===== Stylist Handler =====
type StylistHandler struct {
	service *services.StylistService
}

func NewStylistHandler(service *services.StylistService) *StylistHandler {
	return &StylistHandler{service: service}
}

func (h *StylistHandler) Create(c *gin.Context) {
	var req dto.CreateStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	stylist, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Stylist created", stylist)
}

func (h *StylistHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	stylists, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, stylists, page, pageSize, total)
}

func (h *StylistHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	stylist, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Stylist not found")
		return
	}
	utils.RespondSuccess(c, "", stylist)
}

func (h *StylistHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	stylist, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Stylist updated", stylist)
}

func (h *StylistHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Stylist deleted", nil)
}

// ===== Customer Handler =====
type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	customer, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Customer created", customer)
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	customers, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, customers, page, pageSize, total)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	customer, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Customer not found")
		return
	}
	utils.RespondSuccess(c, "", customer)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	customer, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Customer updated", customer)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Customer deleted", nil)
}

// Unused import suppressor
var _ = models.RoleAdmin
