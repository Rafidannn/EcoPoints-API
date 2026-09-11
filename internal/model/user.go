package model

import "time"

type User struct {
	ID              uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"column:name;not null" json:"name"`
	Email           string     `gorm:"column:email;unique;not null" json:"email"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at" json:"email_verified_at,omitempty"`
	Password        string     `gorm:"column:password;not null" json:"-"`
	Role            string     `gorm:"column:role;default:user;not null" json:"role"`
	PointsBalance   uint64     `gorm:"column:points_balance;default:0;not null" json:"points_balance"`
	RememberToken   *string    `gorm:"column:remember_token" json:"-"`
	CreatedAt       *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
