package service

import (
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
	rewardRepo   repository.RewardRepository
	notifService NotificationService
}

func NewRewardService(rewardRepo repository.RewardRepository, notifService NotificationService) RewardService {
	return &rewardService{
		rewardRepo:   rewardRepo,
		notifService: notifService,
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
		Category:    req.Category,
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
	if req.Category != nil {
		reward.Category = *req.Category
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
		Category:    rew.Category,
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

	redemption.Reward = reward
	res := toRedemptionResponse(redemption)
	return &res, nil
}

func toRedemptionResponse(r *model.RewardRedemption) dto.RedemptionResponse {
	rewardName := ""
	var rewRes *dto.RewardResponse
	if r.Reward != nil {
		rewardName = r.Reward.Name
		res := toRewardResponse(r.Reward)
		rewRes = &res
	}
	userName := ""
	userEmail := ""
	var userRes *dto.UserResponse
	if r.User != nil {
		userName = r.User.Name
		userEmail = r.User.Email
		userRes = &dto.UserResponse{
			ID:             r.User.ID,
			Name:           r.User.Name,
			Email:          r.User.Email,
			Role:           r.User.Role,
			AssignmentArea: r.User.AssignmentArea,
			Address:        r.User.Address,
			WhatsappPhone:  r.User.WhatsappPhone,
			PointsBalance:  r.User.PointsBalance,
			CreatedAt:      r.User.CreatedAt,
			UpdatedAt:      r.User.UpdatedAt,
		}
	}
	return dto.RedemptionResponse{
		ID:          r.ID,
		UserID:      r.UserID,
		UserName:    userName,
		UserEmail:   userEmail,
		RewardID:    r.RewardID,
		RewardName:  rewardName,
		PointsUsed:  r.PointsUsed,
		Status:      r.Status,
		Notes:       r.Notes,
		VoucherCode: r.VoucherCode,
		CreatedAt:   r.CreatedAt,
		User:        userRes,
		Reward:      rewRes,
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
	redemption, err := s.rewardRepo.CompleteRedemption(id, notes)
	if err != nil {
		return nil, err
	}

	rewardName := ""
	if redemption.Reward != nil {
		rewardName = redemption.Reward.Name
	}
	msg := fmt.Sprintf("Penukaran %s telah disetujui!", rewardName)
	if rewardName == "" {
		msg = "Penukaran hadiah kamu telah disetujui!"
	}
	go s.notifService.SendToUser(
		redemption.UserID,
		"Penukaran Hadiah Disetujui 🎁",
		msg,
	)

	res := toRedemptionResponse(redemption)
	return &res, nil
}

func (s *rewardService) RejectRedemption(id uint64, notes *string) (*dto.RedemptionResponse, error) {
	redemption, err := s.rewardRepo.RejectRedemption(id, notes)
	if err != nil {
		return nil, err
	}

	rewardName := ""
	if redemption.Reward != nil {
		rewardName = redemption.Reward.Name
	}
	go s.notifService.SendToUser(
		redemption.UserID,
		"Penukaran Hadiah Ditolak ❌",
		fmt.Sprintf("Penukaran %s ditolak. Poin dan stok hadiah telah dikembalikan.", rewardName),
	)

	res := toRedemptionResponse(redemption)
	return &res, nil
}
