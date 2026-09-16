package dto

import "time"

// WasteDepositItemRequest represents an individual waste item in deposit request
type WasteDepositItemRequest struct {
	WasteTypeID uint64  `json:"waste_type_id"`
	WeightKg    float64 `json:"weight_kg"`
}

// CreateWasteDepositRequest represents payload for submitting a new waste deposit with multiple items
type CreateWasteDepositRequest struct {
	Items       []WasteDepositItemRequest `json:"items"`
	WasteTypeID *uint64                   `json:"waste_type_id"`
	WeightKg    *float64                  `json:"weight_kg"`
	DropPointID *uint64                   `json:"drop_point_id"`
	Notes       *string                   `json:"notes"`
	Photo       *string                   `json:"photo"`
}

// VerifyWasteDepositItemRequest represents verification for an individual item
type VerifyWasteDepositItemRequest struct {
	ItemID   uint64   `json:"item_id"`
	WeightKg *float64 `json:"weight_kg"` // If null or <= 0, uses original submitted weight
}

// VerifyWasteDepositRequest represents payload for a staff/admin verifying a deposit
type VerifyWasteDepositRequest struct {
	Items    []VerifyWasteDepositItemRequest `json:"items"` // optional per item weights
	Status   *string                         `json:"status"`
	WeightKg *float64                        `json:"weight_kg"` // for single item / legacy
	Notes    *string                         `json:"notes"`
}

// WasteDepositItemResponse represents an individual item in deposit response
type WasteDepositItemResponse struct {
	ID               uint64   `json:"id"`
	WasteTypeID      uint64   `json:"waste_type_id"`
	WasteTypeName    string   `json:"waste_type_name"`
	PointsPerKg      uint     `json:"points_per_kg"`
	WeightKg         float64  `json:"weight_kg"`
	OriginalWeightKg float64  `json:"original_weight_kg"`
	ActualWeightKg   *float64 `json:"actual_weight_kg,omitempty"`
	EstimatedPoints  uint     `json:"estimated_points"`
	EarnedPoints     *uint    `json:"earned_points,omitempty"`
}

// WasteDepositResponse represents data returned to client
type WasteDepositResponse struct {
	ID              uint64                     `json:"id"`
	Code            string                     `json:"code"`
	UserID          uint64                     `json:"user_id"`
	UserName        string                     `json:"user_name"`
	Items           []WasteDepositItemResponse `json:"items"`
	TotalWeightKg   float64                    `json:"total_weight_kg"`
	EstimatedPoints uint                       `json:"estimated_points"`
	EarnedPoints    *uint                      `json:"earned_points,omitempty"`
	DropPointID     *uint64                    `json:"drop_point_id,omitempty"`
	DropPointName   *string                    `json:"drop_point_name,omitempty"`
	Status          string                     `json:"status"` // pending, verified, rejected, cancelled
	VerifiedBy      *uint64                    `json:"verified_by,omitempty"`
	VerifierName    *string                    `json:"verifier_name,omitempty"`
	VerifiedAt      *time.Time                 `json:"verified_at,omitempty"`
	Notes           *string                    `json:"notes,omitempty"`
	Photo           *string                    `json:"photo,omitempty"`
	CreatedAt       *time.Time                 `json:"created_at,omitempty"`
}

