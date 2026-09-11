package model

import "time"

type Reward struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"column:name;not null" json:"name"`
	Description *string    `gorm:"column:description;type:text" json:"description,omitempty"`
	PointCost   uint       `gorm:"column:point_cost;not null" json:"point_cost"`
	Stock       int        `gorm:"column:stock;default:0;not null" json:"stock"`
	Image       *string    `gorm:"column:image" json:"image,omitempty"`
	IsActive    bool       `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Reward) TableName() string {
	return "rewards"
}
