package model

import "time"

type WasteDeposit struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint64     `gorm:"column:user_id;not null" json:"user_id"`
	DropPointID *uint64    `gorm:"column:drop_point_id" json:"drop_point_id,omitempty"`
	WasteTypeID uint64     `gorm:"column:waste_type_id;not null" json:"waste_type_id"`
	WeightKg    float64    `gorm:"column:weight_kg;type:decimal(8,2);not null" json:"weight_kg"`
	Photo       *string    `gorm:"column:photo" json:"photo,omitempty"`
	Status      string     `gorm:"column:status;default:pending;not null" json:"status"`
	VerifiedBy  *uint64    `gorm:"column:verified_by" json:"verified_by,omitempty"`
	VerifiedAt  *time.Time `gorm:"column:verified_at" json:"verified_at,omitempty"`
	Notes       *string    `gorm:"column:notes" json:"notes,omitempty"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`

	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DropPoint *DropPoint `gorm:"foreignKey:DropPointID" json:"drop_point,omitempty"`
	WasteType *WasteType `gorm:"foreignKey:WasteTypeID" json:"waste_type,omitempty"`
	Verifier  *User      `gorm:"foreignKey:VerifiedBy" json:"verifier,omitempty"`
}

func (WasteDeposit) TableName() string {
	return "waste_deposits"
}
