package handler

import (
	"net/http"
	"strconv"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	leaderboardService service.LeaderboardService
}

func NewLeaderboardHandler(leaderboardService service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{leaderboardService: leaderboardService}
}

// GetLeaderboard godoc
// @Summary      Get leaderboard
// @Description  Returns top users ranked by points (all-time) or by kg deposited (monthly).
// @Tags         Leaderboard
// @Produce      json
// @Param        period  query  string  false  "Period: 'monthly' or 'all' (default: all)"
// @Param        limit   query  int     false  "Number of entries to return (default: 10, max: 50)"
// @Success      200  {object}  dto.APIResponse{data=[]dto.LeaderboardEntry}
// @Failure      500  {object}  dto.APIErrorResponse
// @Router       /api/v1/leaderboard [get]
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	period := c.DefaultQuery("period", "all")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	entries, err := h.leaderboardService.GetLeaderboard(period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve leaderboard",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Leaderboard retrieved successfully", entries))
}
