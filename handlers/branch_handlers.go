package handlers

import (
	"net/http"

	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/services"
	"hairhaus-pos-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ===== Branch Product Handler =====
type BranchProductHandler struct {
	service *services.BranchProductService
}

func NewBranchProductHandler(service *services.BranchProductService) *BranchProductHandler {
	return &BranchProductHandler{service: service}
}

func (h *BranchProductHandler) Create(c *gin.Context) {
	var req dto.CreateBranchProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	bp, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Branch product created", bp)
}

func (h *BranchProductHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	bps, total, err := h.service.GetByBranch(branchID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, bps, page, pageSize, total)
}

func (h *BranchProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	bp, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Branch product not found")
		return
	}
	utils.RespondSuccess(c, "", bp)
}

func (h *BranchProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateBranchProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	bp, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Branch product updated", bp)
}

func (h *BranchProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Branch product deleted", nil)
}

// ===== Branch Stylist Handler =====
type BranchStylistHandler struct {
	service *services.BranchStylistService
}

func NewBranchStylistHandler(service *services.BranchStylistService) *BranchStylistHandler {
	return &BranchStylistHandler{service: service}
}

func (h *BranchStylistHandler) Create(c *gin.Context) {
	var req dto.CreateBranchStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	bs, err := h.service.Create(req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondCreated(c, "Branch stylist created", bs)
}

func (h *BranchStylistHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid branch ID")
		return
	}
	page, pageSize := utils.GetPaginationParams(c)
	bss, total, err := h.service.GetByBranch(branchID, utils.GetOffset(page, pageSize), pageSize)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondPaginated(c, bss, page, pageSize, total)
}

func (h *BranchStylistHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	bs, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondNotFound(c, "Branch stylist not found")
		return
	}
	utils.RespondSuccess(c, "", bs)
}

func (h *BranchStylistHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	var req dto.UpdateBranchStylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondValidationError(c, err.Error())
		return
	}
	bs, err := h.service.Update(id, req)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Branch stylist updated", bs)
}

func (h *BranchStylistHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondValidationError(c, "Invalid ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondSuccess(c, "Branch stylist deleted", nil)
}
