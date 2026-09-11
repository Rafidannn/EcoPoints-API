package model

import "time"

type RewardRedemption struct {
	ID         uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64     `gorm:"column:user_id;not null" json:"user_id"`
	RewardID   uint64     `gorm:"column:reward_id;not null" json:"reward_id"`
	PointsUsed uint       `gorm:"column:points_used;not null" json:"points_used"`
	Status     string     `gorm:"column:status;default:pending;not null" json:"status"`
	Notes      *string    `gorm:"column:notes;type:text" json:"notes,omitempty"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reward *Reward `gorm:"foreignKey:RewardID" json:"reward,omitempty"`
}

func (RewardRedemption) TableName() string {
	return "reward_redemptions"
}
