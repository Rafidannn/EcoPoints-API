package service

import (
	"errors"
	"time"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"
	"ecopoints-go-api/internal/repository"
)

var (
	ErrDropPointNotFound = errors.New("drop point not found")
)

type DropPointService interface {
	GetAllDropPoints() ([]dto.DropPointResponse, error)
	GetDropPointByID(id uint64) (*dto.DropPointResponse, error)
	CreateDropPoint(req *dto.CreateDropPointRequest) (*dto.DropPointResponse, error)
	UpdateDropPoint(id uint64, req *dto.UpdateDropPointRequest) (*dto.DropPointResponse, error)
	DeleteDropPoint(id uint64) error
}

type dropPointService struct {
	dropPointRepo repository.DropPointRepository
}

func NewDropPointService(dropPointRepo repository.DropPointRepository) DropPointService {
	return &dropPointService{
		dropPointRepo: dropPointRepo,
	}
}

func (s *dropPointService) GetAllDropPoints() ([]dto.DropPointResponse, error) {
	dropPoints, err := s.dropPointRepo.FindAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.DropPointResponse, len(dropPoints))
	for i, dp := range dropPoints {
		response[i] = toDropPointResponse(&dp)
	}

	return response, nil
}

func (s *dropPointService) GetDropPointByID(id uint64) (*dto.DropPointResponse, error) {
	dropPoint, err := s.dropPointRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if dropPoint == nil {
		return nil, ErrDropPointNotFound
	}

	res := toDropPointResponse(dropPoint)
	return &res, nil
}

func (s *dropPointService) CreateDropPoint(req *dto.CreateDropPointRequest) (*dto.DropPointResponse, error) {
	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	dropPoint := &model.DropPoint{
		Name:      req.Name,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		IsActive:  isActive,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := s.dropPointRepo.Create(dropPoint); err != nil {
		return nil, err
	}

	res := toDropPointResponse(dropPoint)
	return &res, nil
}

func (s *dropPointService) UpdateDropPoint(id uint64, req *dto.UpdateDropPointRequest) (*dto.DropPointResponse, error) {
	dropPoint, err := s.dropPointRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if dropPoint == nil {
		return nil, ErrDropPointNotFound
	}

	if req.Name != nil {
		dropPoint.Name = *req.Name
	}
	if req.Address != nil {
		dropPoint.Address = *req.Address
	}
	if req.Latitude != nil {
		dropPoint.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		dropPoint.Longitude = req.Longitude
	}
	if req.IsActive != nil {
		dropPoint.IsActive = *req.IsActive
	}

	now := time.Now()
	dropPoint.UpdatedAt = &now

	if err := s.dropPointRepo.Update(dropPoint); err != nil {
		return nil, err
	}

	res := toDropPointResponse(dropPoint)
	return &res, nil
}

func (s *dropPointService) DeleteDropPoint(id uint64) error {
	dropPoint, err := s.dropPointRepo.FindByID(id)
	if err != nil {
		return err
	}
	if dropPoint == nil {
		return ErrDropPointNotFound
	}

	return s.dropPointRepo.Delete(id)
}

func toDropPointResponse(dp *model.DropPoint) dto.DropPointResponse {
	return dto.DropPointResponse{
		ID:        dp.ID,
		Name:      dp.Name,
		Address:   dp.Address,
		Latitude:  dp.Latitude,
		Longitude: dp.Longitude,
		IsActive:  dp.IsActive,
		CreatedAt: dp.CreatedAt,
		UpdatedAt: dp.UpdatedAt,
	}
}
