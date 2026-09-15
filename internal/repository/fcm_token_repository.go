package repository

import (
	"ecopoints-go-api/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FCMTokenRepository interface {
	Upsert(userID uint64, token, platform string) error
	DeleteByToken(userID uint64, token string) error
	GetTokensByUserID(userID uint64) ([]string, error)
}

type fcmTokenRepository struct {
	db *gorm.DB
}

func NewFCMTokenRepository(db *gorm.DB) FCMTokenRepository {
	return &fcmTokenRepository{db: db}
}

func (r *fcmTokenRepository) Upsert(userID uint64, token, platform string) error {
	now := time.Now()
	t := model.FCMToken{
		UserID:    userID,
		Token:     token,
		Platform:  platform,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	return r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "token"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"user_id":    userID,
				"platform":   platform,
				"updated_at": now,
			}),
		}).
		Create(&t).Error
}

func (r *fcmTokenRepository) DeleteByToken(userID uint64, token string) error {
	return r.db.
		Where("user_id = ? AND token = ?", userID, token).
		Delete(&model.FCMToken{}).Error
}

func (r *fcmTokenRepository) GetTokensByUserID(userID uint64) ([]string, error) {
	var tokens []string
	err := r.db.
		Model(&model.FCMToken{}).
		Where("user_id = ?", userID).
		Pluck("token", &tokens).Error
	return tokens, err
}
