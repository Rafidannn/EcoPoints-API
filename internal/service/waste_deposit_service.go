package service

import (
	"errors"
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
	Reject(depositID uint64, actorID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error)
	Cancel(depositID uint64, userID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error)
}

type wasteDepositService struct {
	repo          repository.WasteDepositRepository
	wasteTypeRepo repository.WasteTypeRepository
	notifService  NotificationService
}

func NewWasteDepositService(repo repository.WasteDepositRepository, wasteTypeRepo repository.WasteTypeRepository, notifService NotificationService) WasteDepositService {
	return &wasteDepositService{
		repo:          repo,
		wasteTypeRepo: wasteTypeRepo,
		notifService:  notifService,
	}
}

func (s *wasteDepositService) toResponse(d *model.WasteDeposit, earnedPoints *uint) dto.WasteDepositResponse {
	code := fmt.Sprintf("ECP-%05d", d.ID)

	userName := ""
	if d.User != nil {
		userName = d.User.Name
	}

	var dropPointName *string
	if d.DropPoint != nil {
		dropPointName = &d.DropPoint.Name
	}

	var verifierName *string
	if d.Verifier != nil {
		verifierName = &d.Verifier.Name
	}

	var totalWeight float64
	var totalEstimatedPoints uint
	var calculatedEarnedPoints uint

	itemsResponse := make([]dto.WasteDepositItemResponse, len(d.Items))
	for i, item := range d.Items {
		wasteTypeName := ""
		pointsPerKg := uint(0)
		if item.WasteType != nil {
			wasteTypeName = item.WasteType.Name
			pointsPerKg = item.WasteType.PointsPerKg
		}

		origWeight := item.OriginalWeightKg
		if origWeight <= 0 {
			origWeight = item.WeightKg
		}

		itemWeight := item.WeightKg
		if item.ActualWeightKg != nil && *item.ActualWeightKg > 0 {
			itemWeight = *item.ActualWeightKg
		}

		estPoints := uint(origWeight * float64(pointsPerKg))
		totalEstimatedPoints += estPoints
		totalWeight += itemWeight

		var itemEarned *uint
		if d.Status == "verified" {
			ep := uint(itemWeight * float64(pointsPerKg))
			itemEarned = &ep
			calculatedEarnedPoints += ep
		}

		itemsResponse[i] = dto.WasteDepositItemResponse{
			ID:               item.ID,
			WasteTypeID:      item.WasteTypeID,
			WasteTypeName:    wasteTypeName,
			PointsPerKg:      pointsPerKg,
			WeightKg:         itemWeight,
			OriginalWeightKg: origWeight,
			ActualWeightKg:   item.ActualWeightKg,
			EstimatedPoints:  estPoints,
			EarnedPoints:     itemEarned,
		}
	}

	var finalEarnedPoints *uint
	if earnedPoints != nil {
		finalEarnedPoints = earnedPoints
	} else if d.Status == "verified" {
		finalEarnedPoints = &calculatedEarnedPoints
	}

	return dto.WasteDepositResponse{
		ID:              d.ID,
		Code:            code,
		UserID:          d.UserID,
		UserName:        userName,
		Items:           itemsResponse,
		TotalWeightKg:   totalWeight,
		EstimatedPoints: totalEstimatedPoints,
		EarnedPoints:    finalEarnedPoints,
		DropPointID:     d.DropPointID,
		DropPointName:   dropPointName,
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
	if len(req.Items) == 0 {
		return nil, errors.New("minimal harus ada 1 jenis sampah yang disetor")
	}

	now := time.Now()
	deposit := model.WasteDeposit{
		UserID:      userID,
		DropPointID: req.DropPointID,
		Photo:       req.Photo,
		Status:      "pending",
		Notes:       req.Notes,
		CreatedAt:   &now,
		UpdatedAt:   &now,
		Items:       make([]model.WasteDepositItem, len(req.Items)),
	}

	for i, itemReq := range req.Items {
		if itemReq.WeightKg <= 0 {
			return nil, fmt.Errorf("berat sampah harus lebih dari 0 kg")
		}

		wasteType, err := s.wasteTypeRepo.FindByID(itemReq.WasteTypeID)
		if err != nil {
			return nil, fmt.Errorf("jenis sampah ID %d tidak valid: %w", itemReq.WasteTypeID, err)
		}

		deposit.Items[i] = model.WasteDepositItem{
			WasteTypeID:      itemReq.WasteTypeID,
			WeightKg:         itemReq.WeightKg,
			OriginalWeightKg: itemReq.WeightKg,
			CreatedAt:        &now,
			UpdatedAt:        &now,
			WasteType:        wasteType,
		}
	}

	if err := s.repo.Create(&deposit); err != nil {
		return nil, fmt.Errorf("gagal membuat setoran sampah: %w", err)
	}

	res := s.toResponse(&deposit, nil)
	return &res, nil
}

func (s *wasteDepositService) GetByID(id uint64) (*dto.WasteDepositResponse, error) {
	deposit, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("setoran tidak ditemukan: %w", err)
	}
	res := s.toResponse(deposit, nil)
	return &res, nil
}

func (s *wasteDepositService) GetMyDeposits(userID uint64) ([]dto.WasteDepositResponse, error) {
	deposits, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	res := make([]dto.WasteDepositResponse, len(deposits))
	for i := range deposits {
		res[i] = s.toResponse(&deposits[i], nil)
	}
	return res, nil
}

func (s *wasteDepositService) GetAll(status string) ([]dto.WasteDepositResponse, error) {
	deposits, err := s.repo.GetAll(status)
	if err != nil {
		return nil, err
	}
	res := make([]dto.WasteDepositResponse, len(deposits))
	for i := range deposits {
		res[i] = s.toResponse(&deposits[i], nil)
	}
	return res, nil
}

func (s *wasteDepositService) Verify(depositID uint64, verifierID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error) {
	itemWeights := make(map[uint64]float64)
	for _, item := range req.Items {
		if item.WeightKg != nil && *item.WeightKg > 0 {
			itemWeights[item.ItemID] = *item.WeightKg
		}
	}

	updated, earnedPoints, err := s.repo.Verify(depositID, verifierID, itemWeights, req.Notes)
	if err != nil {
		return nil, err
	}

	go s.notifService.SendToUser(
		updated.UserID,
		"Setoran Terverifikasi ✅",
		fmt.Sprintf("Setoran #ECP-%05d telah diverifikasi. %d poin berhasil ditambahkan ke akun kamu!", updated.ID, earnedPoints),
	)

	res := s.toResponse(updated, &earnedPoints)
	return &res, nil
}

func (s *wasteDepositService) Reject(depositID uint64, actorID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error) {
	updated, err := s.repo.Reject(depositID, actorID, req.Notes)
	if err != nil {
		return nil, err
	}

	go s.notifService.SendToUser(
		updated.UserID,
		"Setoran Ditolak ❌",
		fmt.Sprintf("Setoran #ECP-%05d ditolak oleh petugas. Silakan hubungi bank sampah untuk informasi lebih lanjut.", updated.ID),
	)

	res := s.toResponse(updated, nil)
	return &res, nil
}

func (s *wasteDepositService) Cancel(depositID uint64, userID uint64, req dto.VerifyWasteDepositRequest) (*dto.WasteDepositResponse, error) {
	updated, err := s.repo.Cancel(depositID, userID, req.Notes)
	if err != nil {
		return nil, err
	}

	res := s.toResponse(updated, nil)
	return &res, nil
}

