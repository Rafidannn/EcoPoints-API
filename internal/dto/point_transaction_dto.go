package dto

import "time"

type PointTransactionResponse struct {
	ID            uint64     `json:"id"`
	UserID        uint64     `json:"user_id"`
	Type          string     `json:"type"`
	Amount        uint       `json:"amount"`
	ReferenceType *string    `json:"reference_type,omitempty"`
	ReferenceID   *uint64    `json:"reference_id,omitempty"`
	Description   *string    `json:"description,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type ReportSummaryResponse struct {
	TotalDeposits       int64                 `json:"total_deposits"`
	TotalWeightKg       float64               `json:"total_weight_kg"`
	TotalPointsIssued   uint64                `json:"total_points_issued"`
	TotalPointsRedeemed uint64                `json:"total_points_redeemed"`
	ByWasteType         []WasteTypeReportItem `json:"by_waste_type"`
	ByDropPoint         []DropPointReportItem `json:"by_drop_point"`
}

type WasteTypeReportItem struct {
	WasteTypeID   uint64  `json:"waste_type_id"`
	WasteTypeName string  `json:"waste_type_name"`
	TotalWeightKg float64 `json:"total_weight_kg"`
	TotalDeposits int     `json:"total_deposits"`
	TotalPoints   uint64  `json:"total_points"`
}

type DropPointReportItem struct {
	DropPointID   uint64  `json:"drop_point_id"`
	DropPointName string  `json:"drop_point_name"`
	TotalWeightKg float64 `json:"total_weight_kg"`
	TotalDeposits int     `json:"total_deposits"`
}
