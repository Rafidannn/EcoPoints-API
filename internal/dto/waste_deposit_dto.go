package dto

import "time"

// CreateWasteDepositRequest represents payload for submitting a new waste deposit
type CreateWasteDepositRequest struct {
	WasteTypeID uint64  `json:"waste_type_id" binding:"required"`
	DropPointID *uint64 `json:"drop_point_id"`
	WeightKg    float64 `json:"weight_kg" binding:"required,gt=0"`
	Notes       *string `json:"notes"`
	Photo       *string `json:"photo"`
}

// VerifyWasteDepositRequest represents payload for a staff/admin verifying a deposit
type VerifyWasteDepositRequest struct {
	WeightKg *float64 `json:"weight_kg"` // If null, uses original submitted weight
	Notes    *string  `json:"notes"`
}

// WasteDepositResponse represents data returned to client
type WasteDepositResponse struct {
	ID               uint64     `json:"id"`
	Code             string     `json:"code"`
	UserID           uint64     `json:"user_id"`
	UserName         string     `json:"user_name"`
	WasteTypeID      uint64     `json:"waste_type_id"`
	WasteTypeName    string     `json:"waste_type_name"`
	PointsPerKg      uint       `json:"points_per_kg"`
	DropPointID      *uint64    `json:"drop_point_id,omitempty"`
	DropPointName    *string    `json:"drop_point_name,omitempty"`
	WeightKg         float64    `json:"weight_kg"`
	OriginalWeightKg float64    `json:"original_weight_kg"`
	ActualWeightKg   *float64   `json:"actual_weight_kg,omitempty"`
	EstimatedPoints  uint       `json:"estimated_points"`
	EarnedPoints     *uint      `json:"earned_points,omitempty"`
	Status           string     `json:"status"` // pending, verified, rejected
	VerifiedBy       *uint64    `json:"verified_by,omitempty"`
	VerifierName     *string    `json:"verifier_name,omitempty"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	Notes            *string    `json:"notes,omitempty"`
	Photo            *string    `json:"photo,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
}
