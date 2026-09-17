package dto

import "time"

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=255" example:"John Doe"`
	Email    string `json:"email" binding:"required,email,max=255" example:"johndoe@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"secret123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@ecopoints.test"`
	Password string `json:"password" binding:"required" example:"password"`
}

type UserResponse struct {
	ID             uint64     `json:"id" example:"1"`
	Name           string     `json:"name" example:"John Doe"`
	Email          string     `json:"email" example:"johndoe@example.com"`
	Role           string     `json:"role" example:"user"`
	PointsBalance  uint64     `json:"points_balance" example:"1500"`
	AssignmentArea *string    `json:"assignment_area,omitempty"`
	Address        *string    `json:"address,omitempty"`
	WhatsappPhone  *string    `json:"whatsapp_phone,omitempty"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token     string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType string       `json:"token_type" example:"Bearer"`
	ExpiresIn int64        `json:"expires_in" example:"259200"`
	User      UserResponse `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"secret123"`
	NewPassword     string `json:"new_password" binding:"required,min=8" example:"newsecret123"`
}

type UpdateProfileRequest struct {
	Name          string  `json:"name" binding:"required,min=2,max=255" example:"John Doe"`
	Email         *string `json:"email" example:"johndoe@example.com"`
	WhatsappPhone *string `json:"whatsapp_phone" example:"081234567890"`
	Address       *string `json:"address" example:"Jl. Sudirman No. 123"`
}

type AdminCreateUserRequest struct {
	Name           string  `json:"name" binding:"required,min=2,max=255"`
	Email          string  `json:"email" binding:"required,email,max=255"`
	Password       string  `json:"password" binding:"required,min=8"`
	Role           string  `json:"role" binding:"required,oneof=user petugas admin"`
	AssignmentArea *string `json:"assignment_area"`
	Address        *string `json:"address"`
	WhatsappPhone  *string `json:"whatsapp_phone"`
}

type AdminUpdateUserRequest struct {
	Name           string  `json:"name" binding:"required,min=2,max=255"`
	Email          string  `json:"email" binding:"required,email,max=255"`
	Password       *string `json:"password"`
	Role           string  `json:"role" binding:"required,oneof=user petugas admin"`
	PointsBalance  uint64  `json:"points_balance"`
	AssignmentArea *string `json:"assignment_area"`
	Address        *string `json:"address"`
	WhatsappPhone  *string `json:"whatsapp_phone"`
}
