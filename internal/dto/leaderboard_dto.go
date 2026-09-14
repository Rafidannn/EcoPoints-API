package dto

// LeaderboardEntry represents a single user in the leaderboard
type LeaderboardEntry struct {
	Rank          int     `json:"rank"`
	UserID        uint64  `json:"user_id"`
	Name          string  `json:"name"`
	PointsBalance uint64  `json:"points_balance"`
	TotalKg       float64 `json:"total_kg"`
}
