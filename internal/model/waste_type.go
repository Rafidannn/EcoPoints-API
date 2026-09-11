package model

import "time"

type WasteType struct {
	ID              uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"column:name;not null" json:"name"`
	UnitPricePerKg  float64    `gorm:"column:unit_price_per_kg;type:decimal(12,2);default:0;not null" json:"unit_price_per_kg"`
	PointsPerKg     uint       `gorm:"column:points_per_kg;default:0;not null" json:"points_per_kg"`
	Description     *string    `gorm:"column:description" json:"description,omitempty"`
	IsActive        bool       `gorm:"column:is_active;default:true;not null" json:"is_active"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (WasteType) TableName() string {
	return "waste_types"
}
