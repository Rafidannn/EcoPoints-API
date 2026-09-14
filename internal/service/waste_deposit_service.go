package service

import (
	"fmt"
	"time"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"
	"ecopoints-go-api/internal/repository"
)

type WasteDepositService interface {
	Create(userID uint64, req dto.CreateWasteDepositRequest) (*dto.WasteDepositResponse, error)
	GetByID(id uint64) (*dto.WasteDepositResponse, error)
	GetMyDeposits(userID uint64) ([]dto.WasteDepositResponse, error)
	GetAll(status string) ([]dto.WasteDepositResponse, error)
	Verify(depositID uint64, verifierID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error)
}

type wasteDepositService struct {
	repo          repository.WasteDepositRepository
	wasteTypeRepo repository.WasteTypeRepository
}

func NewWasteDepositService(repo repository.WasteDepositRepository, wasteTypeRepo repository.WasteTypeRepository) WasteDepositService {
	return &wasteDepositService{
		repo:          repo,
		wasteTypeRepo: wasteTypeRepo,
	}
}

func (s *wasteDepositService) toResponse(d *model.WasteDeposit, earnedPoints *uint) dto.WasteDepositResponse {
	code := fmt.Sprintf("ECP-%05d", d.ID)

	userName := ""
	if d.User != nil {
		userName = d.User.Name
	}

	wasteTypeName := ""
	pointsPerKg := uint(0)
	if d.WasteType != nil {
		wasteTypeName = d.WasteType.Name
		pointsPerKg = d.WasteType.PointsPerKg
	}

	var dropPointName *string
	if d.DropPoint != nil {
		dropPointName = &d.DropPoint.Name
	}

	var verifierName *string
	if d.Verifier != nil {
		verifierName = &d.Verifier.Name
	}

	estimatedPoints := uint(d.WeightKg * float64(pointsPerKg))

	return dto.WasteDepositResponse{
		ID:              d.ID,
		Code:            code,
		UserID:          d.UserID,
		UserName:        userName,
		WasteTypeID:     d.WasteTypeID,
		WasteTypeName:   wasteTypeName,
		PointsPerKg:     pointsPerKg,
		DropPointID:     d.DropPointID,
		DropPointName:   dropPointName,
		WeightKg:        d.WeightKg,
		EstimatedPoints: estimatedPoints,
		EarnedPoints:    earnedPoints,
		Status:          d.Status,
		VerifiedBy:      d.VerifiedBy,
		VerifierName:    verifierName,
		VerifiedAt:      d.VerifiedAt,
		Notes:           d.Notes,
		Photo:           d.Photo,
		CreatedAt:       d.CreatedAt,
	}
}

func (s *wasteDepositService) Create(userID uint64, req dto.CreateWasteDepositRequest) (*dto.WasteDepositResponse, error) {
	// 1. Check waste type exists
	wasteType, err := s.wasteTypeRepo.FindByID(req.WasteTypeID)
	if err != nil {
		return nil, fmt.Errorf("jenis sampah tidak valid: %w", err)
	}

	now := time.Now()
	deposit := model.WasteDeposit{
		UserID:      userID,
		DropPointID: req.DropPointID,
		WasteTypeID: req.WasteTypeID,
		WeightKg:    req.WeightKg,
		Photo:       req.Photo,
		Status:      "pending",
		Notes:       req.Notes,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	if err := s.repo.Create(&deposit); err != nil {
		return nil, fmt.Errorf("gagal membuat setoran sampah: %w", err)
	}

	// Attach preloaded waste type for DTO formatting
	deposit.WasteType = wasteType

	res := s.toResponse(&deposit, nil)
	return &res, nil
}

func (s *wasteDepositService) GetByID(id uint64) (*dto.WasteDepositResponse, error) {
	deposit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("setoran tidak ditemukan: %w", err)
	}
	var earned *uint
	if deposit.Status == "verified" && deposit.WasteType != nil {
		p := uint(deposit.WeightKg * float64(deposit.WasteType.PointsPerKg))
		earned = &p
	}
	res := s.toResponse(deposit, earned)
	return &res, nil
}

func (s *wasteDepositService) GetMyDeposits(userID uint64) ([]dto.WasteDepositResponse, error) {
	deposits, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.WasteDepositResponse, len(deposits))
	for i, d := range deposits {
		var earned *uint
		if d.Status == "verified" && d.WasteType != nil {
			p := uint(d.WeightKg * float64(d.WasteType.PointsPerKg))
			earned = &p
		}
		res[i] = s.toResponse(&d, earned)
	}
	return res, nil
}

func (s *wasteDepositService) GetAll(status string) ([]dto.WasteDepositResponse, error) {
	deposits, err := s.repo.GetAll(status)
	if err != nil {
		return nil, err
	}
	res := make([]dto.WasteDepositResponse, len(deposits))
	for i, d := range deposits {
		var earned *uint
		if d.Status == "verified" && d.WasteType != nil {
			p := uint(d.WeightKg * float64(d.WasteType.PointsPerKg))
			earned = &p
		}
		res[i] = s.toResponse(&d, earned)
	}
	return res, nil
}

func (s *wasteDepositService) Verify(depositID uint64, verifierID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error) {
	actualWeight := 0.0
	if req.WeightKg != nil && *req.WeightKg > 0 {
		actualWeight = *req.WeightKg
	}

	updated, earnedPoints, err := s.repo.Verify(depositID, verifierID, actualWeight, req.Notes)
	if err != nil {
		return nil, err
	}

	res := s.toResponse(updated, &earnedPoints)
	return &res, nil
}
