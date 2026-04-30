package handlers

import (
	"hairhaus-pos-be/clients"
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/services"
	"hairhaus-pos-be/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	service       *services.CustomerService
	loyaltyClient *clients.LoyaltyClient
}

func NewCustomerHandler(service *services.CustomerService, loyaltyClient *clients.LoyaltyClient) *CustomerHandler {
	return &CustomerHandler{
		service:       service,
		loyaltyClient: loyaltyClient,
	}
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
	name := c.Query("name")
	phone := c.Query("phone")
	customers, total, err := h.service.GetAll(utils.GetOffset(page, pageSize), pageSize, name, phone)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, customers, page, pageSize, total)
}

// GetAllSoftDelete godoc
// @Summary List all soft-deleted customers
// @Description Returns a paginated list of customers that have been soft-deleted
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Router /customers/deleted [get]
func (h *CustomerHandler) GetAllSoftDelete(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	customers, total, err := h.service.GetAllDeleted(utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]dto.DeletedCustomerResponse, len(customers))
	for i, cust := range customers {
		response[i] = dto.DeletedCustomerResponse{
			ID:        cust.ID,
			Name:      cust.Name,
			Phone:     cust.Phone,
			DeletedAt: cust.DeletedAt.Time,
		}
	}

	utils.RespondPaginated(c, response, page, pageSize, total)
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

// Delete godoc
// @Summary Delete customer
// @Description Soft delete by default, hard delete if hard_delete=true
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param query query dto.DeleteCustomerRequest false "Delete options"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [delete]
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}

	var query dto.DeleteCustomerRequest
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	hardDelete := query.HardDelete
	var deleteErr error
	if hardDelete {
		deleteErr = h.service.HardDelete(id)
	} else {
		deleteErr = h.service.Delete(id)
	}

	if deleteErr != nil {
		utils.RespondError(c, http.StatusInternalServerError, deleteErr.Error())
		return
	}

	message := "Customer deleted"
	if hardDelete {
		message = "Customer permanently deleted"
	}
	utils.RespondSuccess(c, message, nil)
}

// Identify godoc
// @Summary Identify customer by phone number
// @Description Search ERP database then Loyalty Service for customer identification
// @Tags customers
// @Accept json
// @Produce json
// @Param request body dto.IdentifyCustomerRequest true "Identify Request"
// @Success 200 {object} map[string]interface{}
// @Router /customers/identify [post]
func (h *CustomerHandler) Identify(c *gin.Context) {
	var req dto.IdentifyCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	customer, loyaltyInfo, err := h.service.Identify(c, req.Phone, h.loyaltyClient)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	message := "Customer identified"
	if customer == nil {
		message = "Customer not found"
	}

	utils.RespondSuccess(c, message, gin.H{
		"customer":     customer,
		"loyalty_info": loyaltyInfo,
	})
}

// RegisterLoyalty godoc
// @Summary Register new customer to loyalty system
// @Tags customers
// @Accept json
// @Produce json
// @Param request body dto.RegisterLoyaltyRequest true "Register Request"
// @Success 200 {object} map[string]interface{}
// @Router /customers/loyalty/register [post]
func (h *CustomerHandler) RegisterLoyalty(c *gin.Context) {
	var req dto.RegisterLoyaltyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	customer, requiresVerification, err := h.service.Register(c, req, h.loyaltyClient)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := dto.RegisterLoyaltyResponse{
		Customer:             customer,
		RequiresVerification: requiresVerification,
	}

	message := "Loyalty registered, but requires verification"
	if !requiresVerification {
		message = "Loyalty already verified"
	}

	utils.RespondSuccess(c, message, response)
}

// RequestLoyaltyOTP godoc
// @Summary Request OTP for loyalty
// @Tags customers
// @Accept json
// @Produce json
// @Param phone query string true "Phone Number"
// @Param user_id query string false "Loyalty User ID"
// @Success 200 {object} map[string]interface{}
// @Router /customers/loyalty/otp [post]
func (h *CustomerHandler) RequestLoyaltyOTP(c *gin.Context) {
	phone := c.Query("phone")
	userID := c.Query("user_id")

	if phone == "" {
		utils.RespondValidationError(c, "phone is required")
		return
	}

	if err := h.service.RequestOTP(c, phone, userID, h.loyaltyClient); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, "OTP sent successfully", nil)
}

// VerifyLoyaltyOTP godoc
// @Summary Verify OTP and create customer
// @Tags customers
// @Accept json
// @Produce json
// @Param request body dto.VerifyLoyaltyOTPRequest true "Verify Request"
// @Success 201 {object} models.Customer
// @Router /customers/loyalty/verify [post]
func (h *CustomerHandler) VerifyLoyaltyOTP(c *gin.Context) {
	var req dto.VerifyLoyaltyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	customer, err := h.service.VerifyOTP(c, req, h.loyaltyClient)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondCreated(c, "Customer verified and created", customer)
}

// CheckInLoyalty godoc
// @Summary Check-in a customer using QR code and amount
// @Tags customers
// @Accept json
// @Produce json
// @Param request body dto.LoyaltyCheckInRequest true "Check-in Request"
// @Success 200 {object} dto.LoyaltyCheckInResponse
// @Router /customers/loyalty/check-in [post]
func (h *CustomerHandler) CheckInLoyalty(c *gin.Context) {
	var req dto.LoyaltyCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}

	resp, err := h.loyaltyClient.CheckIn(c, req.Code, req.Amount, req.Notes, req.Metadata)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, "Check-in successful", resp)
}
