package handler

import (
	"errors"
	"net/http"
	"strconv"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type DropPointHandler struct {
	dropPointService service.DropPointService
}

func NewDropPointHandler(dropPointService service.DropPointService) *DropPointHandler {
	return &DropPointHandler{
		dropPointService: dropPointService,
	}
}

// GetAll godoc
// @Summary Get all drop points
// @Description Retrieve list of all drop point locations
// @Tags Drop Points
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.DropPointResponse} "Drop points retrieved successfully"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/drop-points [get]
func (h *DropPointHandler) GetAll(c *gin.Context) {
	dropPoints, err := h.dropPointService.GetAllDropPoints()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve drop points",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Drop points retrieved successfully",
		dropPoints,
	))
}

// GetByID godoc
// @Summary Get drop point by ID
// @Description Retrieve details of a specific drop point location
// @Tags Drop Points
// @Produce json
// @Param id path int true "Drop Point ID"
// @Success 200 {object} dto.APIResponse{data=dto.DropPointResponse} "Drop point retrieved"
// @Failure 400 {object} dto.APIErrorResponse "Invalid ID"
// @Failure 404 {object} dto.APIErrorResponse "Drop point not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/drop-points/{id} [get]
func (h *DropPointHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	dropPoint, err := h.dropPointService.GetDropPointByID(id)
	if err != nil {
		if errors.Is(err, service.ErrDropPointNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"drop_point": "Drop point not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve drop point",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Drop point retrieved successfully",
		dropPoint,
	))
}

// Create godoc
// @Summary Create a new drop point
// @Description Add a new drop point location (Admin only)
// @Tags Drop Points
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Accept json
// @Produce json
// @Param request body dto.CreateDropPointRequest true "Drop Point Data"
// @Success 201 {object} dto.APIResponse{data=dto.DropPointResponse} "Drop point created"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/drop-points [post]
func (h *DropPointHandler) Create(c *gin.Context) {
	var req dto.CreateDropPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	dropPoint, err := h.dropPointService.CreateDropPoint(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to create drop point",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(
		"Drop point created successfully",
		dropPoint,
	))
}

// Update godoc
// @Summary Update an existing drop point
// @Description Update drop point details by ID (Admin only)
// @Tags Drop Points
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Accept json
// @Produce json
// @Param id path int true "Drop Point ID"
// @Param request body dto.UpdateDropPointRequest true "Updated Drop Point Data"
// @Success 200 {object} dto.APIResponse{data=dto.DropPointResponse} "Drop point updated"
// @Failure 400 {object} dto.APIErrorResponse "Invalid input"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "Drop point not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/drop-points/{id} [put]
func (h *DropPointHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	var req dto.UpdateDropPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	dropPoint, err := h.dropPointService.UpdateDropPoint(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrDropPointNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"drop_point": "Drop point not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to update drop point",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Drop point updated successfully",
		dropPoint,
	))
}

// Delete godoc
// @Summary Delete a drop point
// @Description Delete drop point by ID (Admin only)
// @Tags Drop Points
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Produce json
// @Param id path int true "Drop Point ID"
// @Success 200 {object} dto.APIResponse "Drop point deleted"
// @Failure 400 {object} dto.APIErrorResponse "Invalid ID"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "Drop point not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/drop-points/{id} [delete]
func (h *DropPointHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	if err := h.dropPointService.DeleteDropPoint(id); err != nil {
		if errors.Is(err, service.ErrDropPointNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"drop_point": "Drop point not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to delete drop point",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Drop point deleted successfully",
		nil,
	))
}
