package handler

import (
	"net/http"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/repository"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	fcmRepo repository.FCMTokenRepository
}

func NewNotificationHandler(fcmRepo repository.FCMTokenRepository) *NotificationHandler {
	return &NotificationHandler{fcmRepo: fcmRepo}
}

type registerTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform"`
}

type deleteTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func (h *NotificationHandler) RegisterToken(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Sesi tidak valid", nil))
		return
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Sesi tidak valid", nil))
		return
	}

	var req registerTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Token wajib diisi", nil))
		return
	}

	platform := req.Platform
	if platform == "" {
		platform = "android"
	}

	if err := h.fcmRepo.Upsert(userID, req.Token, platform); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menyimpan token", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Token berhasil didaftarkan", nil))
}

func (h *NotificationHandler) DeleteToken(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Sesi tidak valid", nil))
		return
	}
	userID, ok := userIDVal.(uint64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Sesi tidak valid", nil))
		return
	}

	var req deleteTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Token wajib diisi", nil))
		return
	}

	if err := h.fcmRepo.DeleteByToken(userID, req.Token); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menghapus token", nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Token berhasil dihapus", nil))
}
