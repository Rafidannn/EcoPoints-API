package service

import (
	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/repository"
)

type LeaderboardService interface {
	GetLeaderboard(period string, limit int) ([]dto.LeaderboardEntry, error)
}

type leaderboardService struct {
	leaderboardRepo repository.LeaderboardRepository
}

func NewLeaderboardService(leaderboardRepo repository.LeaderboardRepository) LeaderboardService {
	return &leaderboardService{leaderboardRepo: leaderboardRepo}
}

func (s *leaderboardService) GetLeaderboard(period string, limit int) ([]dto.LeaderboardEntry, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	if period == "monthly" {
		return s.leaderboardRepo.GetMonthly(limit)
	}
	return s.leaderboardRepo.GetAllTime(limit)
}
