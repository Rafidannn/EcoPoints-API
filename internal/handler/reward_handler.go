package handler

import (
	"errors"
	"net/http"
	"strconv"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	rewardService service.RewardService
}

func NewRewardHandler(rewardService service.RewardService) *RewardHandler {
	return &RewardHandler{
		rewardService: rewardService,
	}
}

// GetAll godoc
// @Summary Get all rewards
// @Description Retrieve list of all redeemable rewards
// @Tags Rewards
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.RewardResponse} "Rewards retrieved successfully"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/rewards [get]
func (h *RewardHandler) GetAll(c *gin.Context) {
	rewards, err := h.rewardService.GetAllRewards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve rewards",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Rewards retrieved successfully",
		rewards,
	))
}

// GetByID godoc
// @Summary Get reward by ID
// @Description Retrieve details of a specific reward by its ID
// @Tags Rewards
// @Produce json
// @Param id path int true "Reward ID"
// @Success 200 {object} dto.APIResponse{data=dto.RewardResponse} "Reward retrieved"
// @Failure 400 {object} dto.APIErrorResponse "Invalid ID"
// @Failure 404 {object} dto.APIErrorResponse "Reward not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/rewards/{id} [get]
func (h *RewardHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	reward, err := h.rewardService.GetRewardByID(id)
	if err != nil {
		if errors.Is(err, service.ErrRewardNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"reward": "Reward not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve reward",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Reward retrieved successfully",
		reward,
	))
}

// Create godoc
// @Summary Create a new reward
// @Description Add a new reward to catalogue (Admin only)
// @Tags Rewards
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Accept json
// @Produce json
// @Param request body dto.CreateRewardRequest true "Reward Data"
// @Success 201 {object} dto.APIResponse{data=dto.RewardResponse} "Reward created"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/rewards [post]
func (h *RewardHandler) Create(c *gin.Context) {
	var req dto.CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	reward, err := h.rewardService.CreateReward(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to create reward",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(
		"Reward created successfully",
		reward,
	))
}

// Update godoc
// @Summary Update an existing reward
// @Description Update reward by ID (Admin only)
// @Tags Rewards
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Accept json
// @Produce json
// @Param id path int true "Reward ID"
// @Param request body dto.UpdateRewardRequest true "Updated Reward Data"
// @Success 200 {object} dto.APIResponse{data=dto.RewardResponse} "Reward updated"
// @Failure 400 {object} dto.APIErrorResponse "Invalid input"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "Reward not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/rewards/{id} [put]
func (h *RewardHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	var req dto.UpdateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	reward, err := h.rewardService.UpdateReward(id, &req)
	if err != nil {
		if errors.Is(err, service.ErrRewardNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"reward": "Reward not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to update reward",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Reward updated successfully",
		reward,
	))
}

// Delete godoc
// @Summary Delete a reward
// @Description Delete reward by ID (Admin only)
// @Tags Rewards
// @Security BearerAuth
// @Param Authorization header string false "Bearer token"
// @Produce json
// @Param id path int true "Reward ID"
// @Success 200 {object} dto.APIResponse "Reward deleted"
// @Failure 400 {object} dto.APIErrorResponse "Invalid ID"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "Reward not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/rewards/{id} [delete]
func (h *RewardHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Invalid ID format",
			gin.H{"id": "Must be a valid integer ID"},
		))
		return
	}

	if err := h.rewardService.DeleteReward(id); err != nil {
		if errors.Is(err, service.ErrRewardNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"reward": "Reward not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to delete reward",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Reward deleted successfully",
		nil,
	))
}
