package handler

import (
	"errors"
	"net/http"
	"strconv"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WasteTypeHandler struct {
	wasteTypeService service.WasteTypeService
}

func NewWasteTypeHandler(wasteTypeService service.WasteTypeService) *WasteTypeHandler {
	return &WasteTypeHandler{
		wasteTypeService: wasteTypeService,
	}
}

// GetAll godoc
// @Summary Get all waste types
// @Description Retrieve list of all waste types from database
// @Tags Waste Types
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.WasteTypeResponse} "Waste types retrieved successfully"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/waste-types [get]
func (h *WasteTypeHandler) GetAll(c *gin.Context) {
	wasteTypes, err := h.wasteTypeService.GetAllWasteTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve waste types",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Waste types retrieved successfully",
		wasteTypes,
	))
}

// GetByID godoc
// @Summary Get waste type by ID
// @Description Retrieve details of a specific waste type by its ID
// @Tags Waste Types
// @Produce json
// @Param id path int true "Waste Type ID"
// @Success 200 {object} dto.APIResponse{data=dto.WasteTypeResponse} "Waste type retrieved"
// @Failure 400 {object} dto.APIErrorResponse "Invalid ID"
// @Failure 404 {object} dto.APIErrorResponse "Waste type not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/waste-types/{id} [get]
func (h *WasteTypeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	wasteType, err := h.wasteTypeService.GetWasteTypeByID(id)
	if err != nil {
		if errors.Is(err, service.ErrWasteTypeNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"waste_type": "Waste type not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve waste type",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Waste type retrieved successfully",
		wasteType,
	))
}

// Create godoc
// @Summary Create a new waste type
// @Description Create a new waste type (Admin only)
// @Tags Waste Types
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Accept json
// @Produce json
// @Param request body dto.CreateWasteTypeRequest true "Waste Type Data"
// @Success 201 {object} dto.APIResponse{data=dto.WasteTypeResponse} "Waste type created"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/waste-types [post]
func (h *WasteTypeHandler) Create(c *gin.Context) {
	var req dto.CreateWasteTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	wasteType, err := h.wasteTypeService.CreateWasteType(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to create waste type",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(
		"Waste type created successfully",
		wasteType,
	))
}

// Update godoc
// @Summary Update an existing waste type
// @Description Update waste type by ID (Admin only)
// @Tags Waste Types
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Accept json
// @Produce json
// @Param id path int true "Waste Type ID"
// @Param request body dto.UpdateWasteTypeRequest true "Updated Waste Type Data"
// @Success 200 {object} dto.APIResponse{data=dto.WasteTypeResponse} "Waste type updated"
// @Failure 400 {object} dto.APIErrorResponse "Invalid input"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "Waste type not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/waste-types/{id} [put]
func (h *WasteTypeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	var req dto.UpdateWasteTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	wasteType, err := h.wasteTypeService.UpdateWasteType(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrWasteTypeNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"waste_type": "Waste type not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to update waste type",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Waste type updated successfully",
		wasteType,
	))
}

// Delete godoc
// @Summary Delete a waste type
// @Description Delete waste type by ID (Admin only)
// @Tags Waste Types
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Produce json
// @Param id path int true "Waste Type ID"
// @Success 200 {object} dto.APIResponse "Waste type deleted"
// @Failure 400 {object} dto.APIErrorResponse "Invalid ID"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "Waste type not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/waste-types/{id} [delete]
func (h *WasteTypeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	if err := h.wasteTypeService.DeleteWasteType(id); err != nil {
		if errors.Is(err, service.ErrWasteTypeNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"waste_type": "Waste type not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to delete waste type",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Waste type deleted successfully",
		nil,
	))
}
