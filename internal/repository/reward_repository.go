package repository

import (
	"errors"

	"ecopoints-go-api/internal/model"

	"gorm.io/gorm"
)

type RewardRepository interface {
	Create(reward *model.Reward) error
	Update(reward *model.Reward) error
	Delete(id uint64) error
	FindAll() ([]model.Reward, error)
	FindByID(id uint64) (*model.Reward, error)
	Redeem(userID uint64, reward *model.Reward, notes *string) (*model.RewardRedemption, error)
	GetMyRedemptions(userID uint64) ([]model.RewardRedemption, error)
}

type rewardRepository struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) RewardRepository {
	return &rewardRepository{db: db}
}

func (r *rewardRepository) Create(reward *model.Reward) error {
	return r.db.Create(reward).Error
}

func (r *rewardRepository) Update(reward *model.Reward) error {
	return r.db.Save(reward).Error
}

func (r *rewardRepository) Delete(id uint64) error {
	result := r.db.Delete(&model.Reward{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *rewardRepository) FindAll() ([]model.Reward, error) {
	var rewards []model.Reward
	err := r.db.Order("id asc").Find(&rewards).Error
	if err != nil {
		return nil, err
	}
	return rewards, nil
}

func (r *rewardRepository) FindByID(id uint64) (*model.Reward, error) {
	var reward model.Reward
	err := r.db.First(&reward, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reward, nil
}

func (r *rewardRepository) Redeem(userID uint64, reward *model.Reward, notes *string) (*model.RewardRedemption, error) {
	var redemption *model.RewardRedemption

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check stock
		if reward.Stock <= 0 {
			return errors.New("stok hadiah habis")
		}

		// 2. Check user points
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		if int(user.PointsBalance) < int(reward.PointCost) {
			return errors.New("saldo poin tidak mencukupi")
		}

		// 3. Deduct points
		if err := tx.Model(&user).Update("points_balance", user.PointsBalance-uint64(reward.PointCost)).Error; err != nil {
			return err
		}

		// 4. Decrease stock
		if err := tx.Model(reward).Update("stock", reward.Stock-1).Error; err != nil {
			return err
		}

		// 5. Create redemption record
		redemption = &model.RewardRedemption{
			UserID:     userID,
			RewardID:   reward.ID,
			PointsUsed: reward.PointCost,
			Status:     "pending",
			Notes:      notes,
		}
		if err := tx.Create(redemption).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return redemption, nil
}

func (r *rewardRepository) GetMyRedemptions(userID uint64) ([]model.RewardRedemption, error) {
	var redemptions []model.RewardRedemption
	err := r.db.Preload("Reward").Where("user_id = ?", userID).Order("created_at DESC").Find(&redemptions).Error
	if err != nil {
		return nil, err
	}
	return redemptions, nil
}
