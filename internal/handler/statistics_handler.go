package handler

import (
	"net/http"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	statisticsRepo repository.StatisticsRepository
}

func NewStatisticsHandler(statisticsRepo repository.StatisticsRepository) *StatisticsHandler {
	return &StatisticsHandler{statisticsRepo: statisticsRepo}
}

func (h *StatisticsHandler) GetPublic(c *gin.Context) {
	statistics, err := h.statisticsRepo.GetPublicStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil statistik EcoPoints", nil))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse("Statistik EcoPoints berhasil diambil", statistics))
}
