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
