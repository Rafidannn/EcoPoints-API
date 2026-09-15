package dto

type PublicStatisticsResponse struct {
	TotalWeightKg     float64 `json:"total_weight_kg"`
	TotalPointsIssued uint64  `json:"total_points_issued"`
	ActiveDropPoints  int64   `json:"active_drop_points"`
}
