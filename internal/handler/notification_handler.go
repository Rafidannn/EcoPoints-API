package handler

import (
	"fmt"
	"net/http"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/repository"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	fcmRepo      repository.FCMTokenRepository
	notifService service.NotificationService
}

func NewNotificationHandler(fcmRepo repository.FCMTokenRepository, notifService service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		fcmRepo:      fcmRepo,
		notifService: notifService,
	}
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

type sendNotificationRequest struct {
	Target string  `json:"target"` // "all" | "user"
	UserID *uint64 `json:"user_id"`
	Title  string  `json:"title" binding:"required"`
	Body   string  `json:"body" binding:"required"`
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req sendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Judul dan pesan notifikasi wajib diisi", nil))
		return
	}

	if req.Target == "user" && (req.UserID == nil || *req.UserID == 0) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Nasabah tujuan harus dipilih jika target adalah user spesifik", nil))
		return
	}

	if req.Target == "user" {
		err := h.notifService.SendToUser(*req.UserID, req.Title, req.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengirim notifikasi: "+err.Error(), nil))
			return
		}
		c.JSON(http.StatusOK, dto.SuccessResponse("Notifikasi berhasil dikirim ke pengguna", gin.H{"target": "user", "user_id": *req.UserID}))
		return
	}

	// Default broadcast to all
	count, err := h.notifService.SendToAll(req.Title, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengirim notifikasi broadcast: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(fmt.Sprintf("Notifikasi broadcast berhasil dikirim ke %d perangkat nasabah", count), gin.H{"target": "all", "sent_count": count}))
}
