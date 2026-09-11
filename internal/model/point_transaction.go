package model

import "time"

type PointTransaction struct {
	ID            uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID        uint64     `gorm:"column:user_id;not null" json:"user_id"`
	Type          string     `gorm:"column:type;not null" json:"type"`
	Amount        uint       `gorm:"column:amount;not null" json:"amount"`
	ReferenceType *string    `gorm:"column:reference_type" json:"reference_type,omitempty"`
	ReferenceID   *uint64    `gorm:"column:reference_id" json:"reference_id,omitempty"`
	Description   *string    `gorm:"column:description;type:text" json:"description,omitempty"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     *time.Time `gorm:"column:updated_at" json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PointTransaction) TableName() string {
	return "point_transactions"
}
