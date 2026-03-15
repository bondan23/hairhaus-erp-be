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

	utils.RespondSuccess(c, "Customer identified", gin.H{
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

	if err := h.service.Register(c, req, h.loyaltyClient); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, "Loyalty registration initiated", nil)
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
