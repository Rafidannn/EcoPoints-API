package repository

import (
	"ecopoints-go-api/internal/dto"

	"gorm.io/gorm"
)

type LeaderboardRepository interface {
	GetAllTime(limit int) ([]dto.LeaderboardEntry, error)
	GetMonthly(limit int) ([]dto.LeaderboardEntry, error)
}

type leaderboardRepository struct {
	db *gorm.DB
}

func NewLeaderboardRepository(db *gorm.DB) LeaderboardRepository {
	return &leaderboardRepository{db: db}
}

// GetAllTime returns top users ranked by points_balance with total verified kg
func (r *leaderboardRepository) GetAllTime(limit int) ([]dto.LeaderboardEntry, error) {
	type row struct {
		UserID        uint64  `gorm:"column:user_id"`
		Name          string  `gorm:"column:name"`
		PointsBalance uint64  `gorm:"column:points_balance"`
		TotalKg       float64 `gorm:"column:total_kg"`
	}

	var rows []row
	err := r.db.Raw(`
		SELECT
			u.id          AS user_id,
			u.name        AS name,
			u.points_balance AS points_balance,
			COALESCE(SUM(wdi.weight_kg), 0) AS total_kg
		FROM users u
		LEFT JOIN waste_deposits wd
			ON wd.user_id = u.id AND wd.status = 'verified'
		LEFT JOIN waste_deposit_items wdi
			ON wdi.waste_deposit_id = wd.id
		WHERE u.role = 'user'
		GROUP BY u.id, u.name, u.points_balance
		ORDER BY u.points_balance DESC
		LIMIT ?
	`, limit).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	entries := make([]dto.LeaderboardEntry, len(rows))
	for i, r := range rows {
		entries[i] = dto.LeaderboardEntry{
			Rank:          i + 1,
			UserID:        r.UserID,
			Name:          r.Name,
			PointsBalance: r.PointsBalance,
			TotalKg:       r.TotalKg,
		}
	}
	return entries, nil
}

// GetMonthly returns top users by kg deposited in the current month
func (r *leaderboardRepository) GetMonthly(limit int) ([]dto.LeaderboardEntry, error) {
	type row struct {
		UserID        uint64  `gorm:"column:user_id"`
		Name          string  `gorm:"column:name"`
		PointsBalance uint64  `gorm:"column:points_balance"`
		TotalKg       float64 `gorm:"column:total_kg"`
	}

	var rows []row
	err := r.db.Raw(`
		SELECT
			u.id          AS user_id,
			u.name        AS name,
			u.points_balance AS points_balance,
			COALESCE(SUM(wdi.weight_kg), 0) AS total_kg
		FROM users u
		INNER JOIN waste_deposits wd
			ON wd.user_id = u.id
			AND wd.status = 'verified'
			AND MONTH(wd.created_at) = MONTH(NOW())
			AND YEAR(wd.created_at)  = YEAR(NOW())
		INNER JOIN waste_deposit_items wdi
			ON wdi.waste_deposit_id = wd.id
		WHERE u.role = 'user'
		GROUP BY u.id, u.name, u.points_balance
		ORDER BY total_kg DESC
		LIMIT ?
	`, limit).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	entries := make([]dto.LeaderboardEntry, len(rows))
	for i, r := range rows {
		entries[i] = dto.LeaderboardEntry{
			Rank:          i + 1,
			UserID:        r.UserID,
			Name:          r.Name,
			PointsBalance: r.PointsBalance,
			TotalKg:       r.TotalKg,
		}
	}
	return entries, nil
}

