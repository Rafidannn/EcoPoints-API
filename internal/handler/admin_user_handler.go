package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"
	"ecopoints-go-api/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserHandler struct {
	userRepo repository.UserRepository
}

func NewAdminUserHandler(userRepo repository.UserRepository) *AdminUserHandler {
	return &AdminUserHandler{userRepo: userRepo}
}

func (h *AdminUserHandler) GetAll(c *gin.Context) {
	users, err := h.userRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data akun", nil))
		return
	}
	responses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		responses = append(responses, adminUserResponse(&users[i]))
	}
	c.JSON(http.StatusOK, dto.SuccessResponse("Data akun berhasil diambil", responses))
}

func (h *AdminUserHandler) Create(c *gin.Context) {
	var req dto.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validasi gagal", gin.H{"error": err.Error()}))
		return
	}
	if exists, err := h.userRepo.IsEmailExists(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal memeriksa email", nil))
		return
	} else if exists {
		c.JSON(http.StatusConflict, dto.ErrorResponse("Email sudah terdaftar", nil))
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengenkripsi password", nil))
		return
	}
	now := time.Now()
	user := &model.User{Name: req.Name, Email: req.Email, Password: string(hashed), Role: req.Role, AssignmentArea: req.AssignmentArea, Address: req.Address, WhatsappPhone: req.WhatsappPhone, CreatedAt: &now, UpdatedAt: &now}
	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal membuat akun", nil))
		return
	}
	c.JSON(http.StatusCreated, dto.SuccessResponse("Akun berhasil dibuat", adminUserResponse(user)))
}

func (h *AdminUserHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("ID akun tidak valid", nil))
		return
	}
	var req dto.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Validasi gagal", gin.H{"error": err.Error()}))
		return
	}
	user, err := h.userRepo.FindByID(id)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Akun tidak ditemukan", nil))
		return
	}
	user.Name, user.Email, user.Role, user.PointsBalance = req.Name, req.Email, req.Role, req.PointsBalance
	user.AssignmentArea, user.Address, user.WhatsappPhone = req.AssignmentArea, req.Address, req.WhatsappPhone
	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		hashed, hashErr := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengenkripsi password", nil))
			return
		}
		user.Password = string(hashed)
	}
	now := time.Now()
	user.UpdatedAt = &now
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal memperbarui akun", nil))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse("Akun berhasil diperbarui", adminUserResponse(user)))
}

func (h *AdminUserHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("ID akun tidak valid", nil))
		return
	}
	actorID, _ := c.Get("user_id")
	if actor, ok := actorID.(uint64); ok && actor == id {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Akun yang sedang digunakan tidak dapat dihapus", nil))
		return
	}
	if err := h.userRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menghapus akun", nil))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse("Akun berhasil dihapus", nil))
}

func parseID(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func adminUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, PointsBalance: user.PointsBalance, AssignmentArea: user.AssignmentArea, Address: user.Address, WhatsappPhone: user.WhatsappPhone, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}
