package service

import (
	"errors"
	"time"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"
	"ecopoints-go-api/internal/repository"
)

var (
	ErrRewardNotFound = errors.New("reward not found")
)

type RewardService interface {
	GetAllRewards() ([]dto.RewardResponse, error)
	GetRewardByID(id uint64) (*dto.RewardResponse, error)
	CreateReward(req *dto.CreateRewardRequest) (*dto.RewardResponse, error)
	UpdateReward(id uint64, req *dto.UpdateRewardRequest) (*dto.RewardResponse, error)
	DeleteReward(id uint64) error
}

type rewardService struct {
	rewardRepo repository.RewardRepository
}

func NewRewardService(rewardRepo repository.RewardRepository) RewardService {
	return &rewardService{
		rewardRepo: rewardRepo,
	}
}

func (s *rewardService) GetAllRewards() ([]dto.RewardResponse, error) {
	rewards, err := s.rewardRepo.FindAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.RewardResponse, len(rewards))
	for i, rew := range rewards {
		response[i] = toRewardResponse(&rew)
	}

	return response, nil
}

func (s *rewardService) GetRewardByID(id uint64) (*dto.RewardResponse, error) {
	reward, err := s.rewardRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, ErrRewardNotFound
	}

	res := toRewardResponse(reward)
	return &res, nil
}

func (s *rewardService) CreateReward(req *dto.CreateRewardRequest) (*dto.RewardResponse, error) {
	now := time.Now()
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	reward := &model.Reward{
		Name:        req.Name,
		Description: req.Description,
		PointCost:   req.PointCost,
		Stock:       req.Stock,
		Image:       req.Image,
		IsActive:    isActive,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	if err := s.rewardRepo.Create(reward); err != nil {
		return nil, err
	}

	res := toRewardResponse(reward)
	return &res, nil
}

func (s *rewardService) UpdateReward(id uint64, req *dto.UpdateRewardRequest) (*dto.RewardResponse, error) {
	reward, err := s.rewardRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, ErrRewardNotFound
	}

	if req.Name != nil {
		reward.Name = *req.Name
	}
	if req.Description != nil {
		reward.Description = req.Description
	}
	if req.PointCost != nil {
		reward.PointCost = *req.PointCost
	}
	if req.Stock != nil {
		reward.Stock = *req.Stock
	}
	if req.Image != nil {
		reward.Image = req.Image
	}
	if req.IsActive != nil {
		reward.IsActive = *req.IsActive
	}

	now := time.Now()
	reward.UpdatedAt = &now

	if err := s.rewardRepo.Update(reward); err != nil {
		return nil, err
	}

	res := toRewardResponse(reward)
	return &res, nil
}

func (s *rewardService) DeleteReward(id uint64) error {
	reward, err := s.rewardRepo.FindByID(id)
	if err != nil {
		return err
	}
	if reward == nil {
		return ErrRewardNotFound
	}

	return s.rewardRepo.Delete(id)
}

func toRewardResponse(rew *model.Reward) dto.RewardResponse {
	return dto.RewardResponse{
		ID:          rew.ID,
		Name:        rew.Name,
		Description: rew.Description,
		PointCost:   rew.PointCost,
		Stock:       rew.Stock,
		Image:       rew.Image,
		IsActive:    rew.IsActive,
		CreatedAt:   rew.CreatedAt,
		UpdatedAt:   rew.UpdatedAt,
	}
}
