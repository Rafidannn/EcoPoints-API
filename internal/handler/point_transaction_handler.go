package handler

import (
	"net/http"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type PointTransactionHandler struct {
	pointTransactionRepo repository.PointTransactionRepository
}

func NewPointTransactionHandler(pointTransactionRepo repository.PointTransactionRepository) *PointTransactionHandler {
	return &PointTransactionHandler{pointTransactionRepo: pointTransactionRepo}
}

func (h *PointTransactionHandler) GetMyTransactions(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Unauthorized", nil))
		return
	}

	userID, ok := userIDVal.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Internal server error", nil))
		return
	}

	transactions, err := h.pointTransactionRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil transaksi poin", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Transaksi poin berhasil diambil", transactions))
}

func (h *PointTransactionHandler) GetAll(c *gin.Context) {
	transactions, err := h.pointTransactionRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil transaksi poin", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Semua transaksi poin berhasil diambil", transactions))
}

func (h *PointTransactionHandler) GetReportSummary(c *gin.Context) {
	summary, err := h.pointTransactionRepo.GetReportSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil laporan", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Laporan EcoPoints berhasil diambil", summary))
}
