package model

import "time"

type DropPoint struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string     `gorm:"column:name;not null" json:"name"`
	Address   string     `gorm:"column:address;type:text;not null" json:"address"`
	Latitude  *float64   `gorm:"column:latitude;type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude *float64   `gorm:"column:longitude;type:decimal(11,8)" json:"longitude,omitempty"`
	IsActive  bool       `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (DropPoint) TableName() string {
	return "drop_points"
}
