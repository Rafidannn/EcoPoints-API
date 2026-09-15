package repository

import (
	"fmt"
	"strings"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"

	"gorm.io/gorm"
)

func NormalizePointType(input string) string {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "credit", "earned":
		return "credit"
	case "debit", "spent":
		return "debit"
	default:
		return strings.ToLower(strings.TrimSpace(input))
	}
}

type PointTransactionRepository interface {
	GetByUserID(userID uint64) ([]model.PointTransaction, error)
	GetAll() ([]model.PointTransaction, error)
	GetReportSummary() (*dto.ReportSummaryResponse, error)
}

type pointTransactionRepository struct {
	db *gorm.DB
}

func NewPointTransactionRepository(db *gorm.DB) PointTransactionRepository {
	return &pointTransactionRepository{db: db}
}

func (r *pointTransactionRepository) GetByUserID(userID uint64) ([]model.PointTransaction, error) {
	var txns []model.PointTransaction
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&txns).Error; err != nil {
		return nil, err
	}
	for i := range txns {
		txns[i].Type = NormalizePointType(txns[i].Type)
	}
	return txns, nil
}

func (r *pointTransactionRepository) GetAll() ([]model.PointTransaction, error) {
	var txns []model.PointTransaction
	if err := r.db.Order("created_at DESC").Find(&txns).Error; err != nil {
		return nil, err
	}
	for i := range txns {
		txns[i].Type = NormalizePointType(txns[i].Type)
	}
	return txns, nil
}

func (r *pointTransactionRepository) GetReportSummary() (*dto.ReportSummaryResponse, error) {
	var summary dto.ReportSummaryResponse

	var totalDeposits int64
	if err := r.db.Table("waste_deposits").Where("status = ?", "verified").Count(&totalDeposits).Error; err != nil {
		return nil, err
	}
	summary.TotalDeposits = totalDeposits

	if err := r.db.Table("waste_deposits").Where("status = ?", "verified").Select("COALESCE(SUM(weight_kg), 0)").Scan(&summary.TotalWeightKg).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("point_transactions").Where("type IN (?, ?)", "credit", "earned").Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalPointsIssued).Error; err != nil {
		return nil, err
	}
	if err := r.db.Table("point_transactions").Where("type IN (?, ?)", "debit", "spent").Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalPointsRedeemed).Error; err != nil {
		return nil, err
	}

	var byWasteType []dto.WasteTypeReportItem
	if err := r.db.Table("waste_deposits wd").
		Select("wt.id AS waste_type_id, wt.name AS waste_type_name, COALESCE(SUM(wd.weight_kg),0) AS total_weight_kg, COUNT(wd.id) AS total_deposits, 0 AS total_points").
		Joins("LEFT JOIN waste_types wt ON wt.id = wd.waste_type_id").
		Where("wd.status = ?", "verified").
		Group("wt.id, wt.name").
		Scan(&byWasteType).Error; err != nil {
		return nil, err
	}
	summary.ByWasteType = byWasteType

	var byDropPoint []dto.DropPointReportItem
	if err := r.db.Table("waste_deposits wd").
		Select("dp.id AS drop_point_id, dp.name AS drop_point_name, COALESCE(SUM(wd.weight_kg),0) AS total_weight_kg, COUNT(wd.id) AS total_deposits").
		Joins("LEFT JOIN drop_points dp ON dp.id = wd.drop_point_id").
		Where("wd.status = ?", "verified").
		Group("dp.id, dp.name").
		Scan(&byDropPoint).Error; err != nil {
		return nil, err
	}
	summary.ByDropPoint = byDropPoint

	return &summary, nil
}

func (r *pointTransactionRepository) FindByUserID(userID uint64) ([]model.PointTransaction, error) {
	return r.GetByUserID(userID)
}

func (r *pointTransactionRepository) ListAll() ([]model.PointTransaction, error) {
	return r.GetAll()
}

func (r *pointTransactionRepository) Summary() (*dto.ReportSummaryResponse, error) {
	return r.GetReportSummary()
}

func (r *pointTransactionRepository) BuildReportSummary() (*dto.ReportSummaryResponse, error) {
	return r.GetReportSummary()
}

func BuildReportSummaryFromTransactions() string {
	return fmt.Sprintf("summary")
}
