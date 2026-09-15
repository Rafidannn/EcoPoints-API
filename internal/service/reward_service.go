package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/model"
	"ecopoints-go-api/internal/repository"
)

var (
	ErrRewardNotFound     = errors.New("reward not found")
	ErrInsufficientPoints = errors.New("saldo poin tidak mencukupi")
	ErrOutOfStock         = errors.New("stok hadiah habis")
)

type RewardService interface {
	GetAllRewards() ([]dto.RewardResponse, error)
	GetRewardByID(id uint64) (*dto.RewardResponse, error)
	CreateReward(req *dto.CreateRewardRequest) (*dto.RewardResponse, error)
	UpdateReward(id uint64, req *dto.UpdateRewardRequest) (*dto.RewardResponse, error)
	DeleteReward(id uint64) error
	RedeemReward(userID uint64, rewardID uint64, notes *string) (*dto.RedemptionResponse, error)
	GetMyRedemptions(userID uint64) ([]dto.RedemptionResponse, error)
	GetAllRedemptions() ([]dto.RedemptionResponse, error)
	CompleteRedemption(id uint64, notes *string) (*dto.RedemptionResponse, error)
	RejectRedemption(id uint64, notes *string) (*dto.RedemptionResponse, error)
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

func (s *rewardService) RedeemReward(userID uint64, rewardID uint64, notes *string) (*dto.RedemptionResponse, error) {
	reward, err := s.rewardRepo.FindByID(rewardID)
	if err != nil {
		return nil, err
	}
	if reward == nil {
		return nil, ErrRewardNotFound
	}
	if !reward.IsActive {
		return nil, ErrOutOfStock
	}

	redemption, err := s.rewardRepo.Redeem(userID, reward, notes)
	if err != nil {
		return nil, err
	}

	rewardName := reward.Name
	res := &dto.RedemptionResponse{
		ID:         redemption.ID,
		UserID:     redemption.UserID,
		UserName:   "",
		RewardID:   redemption.RewardID,
		RewardName: rewardName,
		PointsUsed: redemption.PointsUsed,
		Status:     redemption.Status,
		Notes:      redemption.Notes,
		CreatedAt:  redemption.CreatedAt,
	}
	return res, nil
}

func toRedemptionResponse(r *model.RewardRedemption) dto.RedemptionResponse {
	rewardName := ""
	if r.Reward != nil {
		rewardName = r.Reward.Name
	}
	userName := ""
	if r.User != nil {
		userName = r.User.Name
	}
	return dto.RedemptionResponse{
		ID:          r.ID,
		UserID:      r.UserID,
		UserName:    userName,
		RewardID:    r.RewardID,
		RewardName:  rewardName,
		PointsUsed:  r.PointsUsed,
		Status:      r.Status,
		Notes:       r.Notes,
		VoucherCode: r.VoucherCode,
		CreatedAt:   r.CreatedAt,
	}
}

func (s *rewardService) GetMyRedemptions(userID uint64) ([]dto.RedemptionResponse, error) {
	redemptions, err := s.rewardRepo.GetMyRedemptions(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.RedemptionResponse, len(redemptions))
	for i, r := range redemptions {
		result[i] = toRedemptionResponse(&r)
	}
	return result, nil
}

func (s *rewardService) GetAllRedemptions() ([]dto.RedemptionResponse, error) {
	redemptions, err := s.rewardRepo.GetAllRedemptions()
	if err != nil {
		return nil, err
	}

	result := make([]dto.RedemptionResponse, len(redemptions))
	for i, r := range redemptions {
		result[i] = toRedemptionResponse(&r)
	}
	return result, nil
}

func (s *rewardService) CompleteRedemption(id uint64, notes *string) (*dto.RedemptionResponse, error) {
	voucherCode, err := generateVoucherCode()
	if err != nil {
		return nil, err
	}

	redemption, err := s.rewardRepo.CompleteRedemption(id, notes, voucherCode)
	if err != nil {
		return nil, err
	}

	res := toRedemptionResponse(redemption)
	return &res, nil
}

func (s *rewardService) RejectRedemption(id uint64, notes *string) (*dto.RedemptionResponse, error) {
	redemption, err := s.rewardRepo.RejectRedemption(id, notes)
	if err != nil {
		return nil, err
	}

	res := toRedemptionResponse(redemption)
	return &res, nil
}

func generateVoucherCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return fmt.Sprintf("EP-RDM-%s", string(b)), nil
}
