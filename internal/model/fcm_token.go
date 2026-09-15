package model

import "time"

type FCMToken struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"not null;index"`
	Token     string     `gorm:"type:varchar(512);uniqueIndex;not null"`
	Platform  string     `gorm:"type:varchar(32);not null;default:android"`
	CreatedAt *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
	User      *User      `gorm:"foreignKey:UserID"`
}
