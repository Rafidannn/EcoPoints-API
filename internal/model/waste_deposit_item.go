package model

import "time"

type WasteDepositItem struct {
	ID               uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	WasteDepositID   uint64     `gorm:"column:waste_deposit_id;not null;index" json:"waste_deposit_id"`
	WasteTypeID      uint64     `gorm:"column:waste_type_id;not null" json:"waste_type_id"`
	WeightKg         float64    `gorm:"column:weight_kg;type:decimal(8,2);not null" json:"weight_kg"`
	OriginalWeightKg float64    `gorm:"column:original_weight_kg;type:decimal(8,2);not null;default:0" json:"original_weight_kg"`
	ActualWeightKg   *float64   `gorm:"column:actual_weight_kg;type:decimal(8,2)" json:"actual_weight_kg,omitempty"`
	CreatedAt        *time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`

	WasteDeposit *WasteDeposit `gorm:"foreignKey:WasteDepositID" json:"waste_deposit,omitempty"`
	WasteType    *WasteType    `gorm:"foreignKey:WasteTypeID" json:"waste_type,omitempty"`
}

func (WasteDepositItem) TableName() string {
	return "waste_deposit_items"
}
