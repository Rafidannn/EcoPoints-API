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
	ID            uint64     `json:"id" example:"1"`
	Name          string     `json:"name" example:"John Doe"`
	Email         string     `json:"email" example:"johndoe@example.com"`
	Role          string     `json:"role" example:"user"`
	PointsBalance uint64     `json:"points_balance" example:"1500"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token     string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType string       `json:"token_type" example:"Bearer"`
	ExpiresIn int64        `json:"expires_in" example:"259200"`
	User      UserResponse `json:"user"`
}
