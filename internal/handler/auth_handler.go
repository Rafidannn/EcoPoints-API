package handler

import (
	"errors"
	"net/http"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with default 'user' role
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User Registration Data"
// @Success 201 {object} dto.APIResponse{data=dto.UserResponse} "User successfully registered"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 409 {object} dto.APIErrorResponse "Email already exists"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	userResponse, err := h.authService.Register(&req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, dto.ErrorResponse(
				err.Error(),
				gin.H{"email": "Email is already in use"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to register user",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(
		"User registered successfully",
		userResponse,
	))
}

// Login godoc
// @Summary Login user
// @Description Authenticate user via email and Laravel bcrypt password, returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} dto.APIResponse{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 401 {object} dto.APIErrorResponse "Invalid credentials"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(
				err.Error(),
				gin.H{"auth": "Invalid email or password"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to login",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Login successful",
		loginResponse,
	))
}

// Me godoc
// @Summary Get current user profile
// @Description Retrieve profile of the currently authenticated user
// @Tags Auth
// @Security BearerAuth
// @Param Authorization header string false "Bearer token (Format: Bearer <token>)"
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse} "User profile retrieved"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 404 {object} dto.APIErrorResponse "User not found"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(
			"Unauthorized",
			gin.H{"auth": "User ID not found in token"},
		))
		return
	}

	userID, ok := userIDVal.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Internal server error",
			gin.H{"auth": "Invalid user ID type"},
		))
		return
	}

	userResponse, err := h.authService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"user": "User record not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Failed to retrieve user profile",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"User profile retrieved successfully",
		userResponse,
	))
}

// UpdateProfile godoc
// @Summary Update current user profile
// @Description Update name, whatsapp phone, address for authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "Update Profile Request"
// @Success 200 {object} dto.APIResponse{data=dto.UserResponse} "User profile updated"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized"
// @Failure 409 {object} dto.APIErrorResponse "Email already exists"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(
			"Unauthorized",
			gin.H{"auth": "User ID not found in token"},
		))
		return
	}

	var userID uint64
	switch v := userIDVal.(type) {
	case uint64:
		userID = v
	case float64:
		userID = uint64(v)
	case int:
		userID = uint64(v)
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Internal server error",
			gin.H{"auth": "Invalid user ID type"},
		))
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	userResponse, err := h.authService.UpdateProfile(userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, dto.ErrorResponse(
				err.Error(),
				gin.H{"email": "Email is already in use"},
			))
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"user": "User record not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Gagal memperbarui profil",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Profil berhasil diperbarui",
		userResponse,
	))
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for the currently authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} dto.APIResponse "Password updated successfully"
// @Failure 400 {object} dto.APIErrorResponse "Validation error"
// @Failure 401 {object} dto.APIErrorResponse "Unauthorized or wrong current password"
// @Failure 500 {object} dto.APIErrorResponse "Internal server error"
// @Router /api/v1/auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse(
			"Unauthorized",
			gin.H{"auth": "User ID not found in token"},
		))
		return
	}

	var userID uint64
	switch v := userIDVal.(type) {
	case uint64:
		userID = v
	case float64:
		userID = uint64(v)
	case int:
		userID = uint64(v)
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Internal server error",
			gin.H{"auth": "Invalid user ID type"},
		))
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			"Validation failed",
			formatValidationError(err),
		))
		return
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrWrongCurrentPassword) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse(
				err.Error(),
				gin.H{"current_password": "Kata sandi saat ini salah"},
			))
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse(
				err.Error(),
				gin.H{"user": "User record not found"},
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse(
			"Gagal mengubah kata sandi",
			gin.H{"error": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(
		"Kata sandi berhasil diubah",
		nil,
	))
}


func formatValidationError(err error) gin.H {
	errMap := gin.H{}
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			errMap[fieldErr.Field()] = fieldErr.Tag()
		}
	} else {
		errMap["general"] = err.Error()
	}
	return errMap
}
