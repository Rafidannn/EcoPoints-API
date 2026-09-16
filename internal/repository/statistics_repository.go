package repository

import (
	"ecopoints-go-api/internal/dto"

	"gorm.io/gorm"
)

type StatisticsRepository interface {
	GetPublicStatistics() (*dto.PublicStatisticsResponse, error)
}

type statisticsRepository struct {
	db *gorm.DB
}

func NewStatisticsRepository(db *gorm.DB) StatisticsRepository {
	return &statisticsRepository{db: db}
}

func (r *statisticsRepository) GetPublicStatistics() (*dto.PublicStatisticsResponse, error) {
	var result dto.PublicStatisticsResponse
	if err := r.db.Table("waste_deposit_items wdi").
		Joins("INNER JOIN waste_deposits wd ON wd.id = wdi.waste_deposit_id").
		Where("wd.status = ?", "verified").
		Select("COALESCE(SUM(wdi.weight_kg), 0)").
		Scan(&result.TotalWeightKg).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("point_transactions").Where("type IN (?, ?)", "credit", "earned").Select("COALESCE(SUM(amount), 0)").Scan(&result.TotalPointsIssued).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("drop_points").Where("is_active = ?", true).Count(&result.ActiveDropPoints).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

