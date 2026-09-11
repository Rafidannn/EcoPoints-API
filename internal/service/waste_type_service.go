package service

import (
	"errors"
	"time"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"
	"ecopoints-go-api/internal/repository"
)

var (
	ErrWasteTypeNotFound = errors.New("waste type not found")
)

type WasteTypeService interface {
	GetAllWasteTypes() ([]dto.WasteTypeResponse, error)
	GetWasteTypeByID(id uint64) (*dto.WasteTypeResponse, error)
	CreateWasteType(req *dto.CreateWasteTypeRequest) (*dto.WasteTypeResponse, error)
	UpdateWasteType(id uint64, req *dto.UpdateWasteTypeRequest) (*dto.WasteTypeResponse, error)
	DeleteWasteType(id uint64) error
}

type wasteTypeService struct {
	wasteTypeRepo repository.WasteTypeRepository
}

func NewWasteTypeService(wasteTypeRepo repository.WasteTypeRepository) WasteTypeService {
	return &wasteTypeService{
		wasteTypeRepo: wasteTypeRepo,
	}
}

func (s *wasteTypeService) GetAllWasteTypes() ([]dto.WasteTypeResponse, error) {
	wasteTypes, err := s.wasteTypeRepo.FindAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.WasteTypeResponse, len(wasteTypes))
	for i, wt := range wasteTypes {
		response[i] = toWasteTypeResponse(&wt)
	}

	return response, nil
}

func (s *wasteTypeService) GetWasteTypeByID(id uint64) (*dto.WasteTypeResponse, error) {
	wasteType, err := s.wasteTypeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if wasteType == nil {
		return nil, ErrWasteTypeNotFound
	}

	res := toWasteTypeResponse(wasteType)
	return &res, nil
}

func (s *wasteTypeService) CreateWasteType(req *dto.CreateWasteTypeRequest) (*dto.WasteTypeResponse, error) {
	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	wasteType := &model.WasteType{
		Name:           req.Name,
		UnitPricePerKg: req.UnitPricePerKg,
		PointsPerKg:    req.PointsPerKg,
		Description:    req.Description,
		IsActive:       isActive,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}

	if err := s.wasteTypeRepo.Create(wasteType); err != nil {
		return nil, err
	}

	res := toWasteTypeResponse(wasteType)
	return &res, nil
}

func (s *wasteTypeService) UpdateWasteType(id uint64, req *dto.UpdateWasteTypeRequest) (*dto.WasteTypeResponse, error) {
	wasteType, err := s.wasteTypeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if wasteType == nil {
		return nil, ErrWasteTypeNotFound
	}

	if req.Name != nil {
		wasteType.Name = *req.Name
	}
	if req.UnitPricePerKg != nil {
		wasteType.UnitPricePerKg = *req.UnitPricePerKg
	}
	if req.PointsPerKg != nil {
		wasteType.PointsPerKg = *req.PointsPerKg
	}
	if req.Description != nil {
		wasteType.Description = req.Description
	}
	if req.IsActive != nil {
		wasteType.IsActive = *req.IsActive
	}

	now := time.Now()
	wasteType.UpdatedAt = &now

	if err := s.wasteTypeRepo.Update(wasteType); err != nil {
		return nil, err
	}

	res := toWasteTypeResponse(wasteType)
	return &res, nil
}

func (s *wasteTypeService) DeleteWasteType(id uint64) error {
	wasteType, err := s.wasteTypeRepo.FindByID(id)
	if err != nil {
		return err
	}
	if wasteType == nil {
		return ErrWasteTypeNotFound
	}

	return s.wasteTypeRepo.Delete(id)
}

func toWasteTypeResponse(wt *model.WasteType) dto.WasteTypeResponse {
	return dto.WasteTypeResponse{
		ID:             wt.ID,
		Name:           wt.Name,
		UnitPricePerKg: wt.UnitPricePerKg,
		PointsPerKg:    wt.PointsPerKg,
		Description:    wt.Description,
		IsActive:       wt.IsActive,
		CreatedAt:      wt.CreatedAt,
		UpdatedAt:      wt.UpdatedAt,
	}
}
