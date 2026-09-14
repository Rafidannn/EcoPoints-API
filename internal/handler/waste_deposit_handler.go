package handler

import (
	"net/http"
	"strconv"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WasteDepositHandler struct {
	wasteDepositService service.WasteDepositService
}

func NewWasteDepositHandler(wasteDepositService service.WasteDepositService) *WasteDepositHandler {
	return &WasteDepositHandler{
		wasteDepositService: wasteDepositService,
	}
}

// Create godoc
// @Summary Submit a new waste deposit
// @Description Creates a new pending waste deposit request
// @Tags Waste Deposits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateWasteDepositRequest true "Waste Deposit Data"
// @Success 201 {object} dto.APIResponse{data=dto.WasteDepositResponse}
// @Failure 400 {object} dto.APIErrorResponse
// @Failure 500 {object} dto.APIErrorResponse
// @Router /api/v1/waste-deposits [post]
func (h *WasteDepositHandler) Create(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Unauthorized", nil))
		return
	}
	userID := userIDVal.(uint64)

	var req dto.CreateWasteDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validasi gagal", gin.H{"error": err.Error()}))
		return
	}

	res, err := h.wasteDepositService.Create(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Formulir setor sampah berhasil dikirim!", res))
}

// GetAll godoc
// @Summary List waste deposits
// @Description Returns waste deposits. Normal users see only their deposits. Admin/Petugas see all.
// @Tags Waste Deposits
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter status: pending, verified"
// @Success 200 {object} dto.APIResponse{data=[]dto.WasteDepositResponse}
// @Failure 500 {object} dto.APIErrorResponse
// @Router /api/v1/waste-deposits [get]
func (h *WasteDepositHandler) GetAll(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Unauthorized", nil))
		return
	}
	userID := userIDVal.(uint64)
	role := c.GetString("role")

	if role == "admin" || role == "petugas" {
		status := c.Query("status")
		res, err := h.wasteDepositService.GetAll(status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data setoran", gin.H{"error": err.Error()}))
			return
		}
		c.JSON(http.StatusOK, dto.SuccessResponse("Data setoran berhasil diambil", res))
		return
	}

	// Regular user gets only their deposits
	res, err := h.wasteDepositService.GetMyDeposits(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil riwayat setoran", gin.H{"error": err.Error()}))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse("Riwayat setoran berhasil diambil", res))
}

// GetByID godoc
// @Summary Get waste deposit details
// @Tags Waste Deposits
// @Produce json
// @Security BearerAuth
// @Param id path int true "Waste Deposit ID"
// @Success 200 {object} dto.APIResponse{data=dto.WasteDepositResponse}
// @Failure 404 {object} dto.APIErrorResponse
// @Router /api/v1/waste-deposits/{id} [get]
func (h *WasteDepositHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("ID tidak valid", nil))
		return
	}

	res, err := h.wasteDepositService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Data setoran tidak ditemukan", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Data setoran berhasil diambil", res))
}

// Verify godoc
// @Summary Verify waste deposit (Staff/Admin only)
// @Description Confirms deposit, assigns actual weight, calculates and awards points to user
// @Tags Waste Deposits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Waste Deposit ID"
// @Param request body dto.VerifyWasteDepositRequest false "Verification Data"
// @Success 200 {object} dto.APIResponse{data=dto.WasteDepositResponse}
// @Failure 400 {object} dto.APIErrorResponse
// @Failure 500 {object} dto.APIErrorResponse
// @Router /api/v1/waste-deposits/{id}/verify [put]
func (h *WasteDepositHandler) Verify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("ID tidak valid", nil))
		return
	}

	verifierIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Unauthorized", nil))
		return
	}
	verifierID := verifierIDVal.(uint64)

	var req dto.VerifyWasteDepositRequest
	_ = c.ShouldBindJSON(&req) // optional body

	res, err := h.wasteDepositService.Verify(id, verifierID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Setoran sampah berhasil diverifikasi & poin telah ditambahkan!", res))
}
