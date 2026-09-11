package dto

import "time"

type CreateRewardRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=255" example:"Voucher Belanja Rp 50.000"`
	Description *string `json:"description" example:"Voucher belanja dapat digunakan di minimarket rekanan"`
	PointCost   uint    `json:"point_cost" binding:"required,gt=0" example:"500"`
	Stock       int     `json:"stock" binding:"required,gte=0" example:"20"`
	Image       *string `json:"image" example:"rewards/voucher50k.png"`
	IsActive    *bool   `json:"is_active" example:"true"`
}

type UpdateRewardRequest struct {
	Name        *string `json:"name,omitempty" example:"Voucher Belanja Rp 50.000 Promo"`
	Description *string `json:"description,omitempty" example:"Voucher belanja berlaku nasional"`
	PointCost   *uint   `json:"point_cost,omitempty" example:"450"`
	Stock       *int    `json:"stock,omitempty" example:"25"`
	Image       *string `json:"image,omitempty" example:"rewards/voucher50k_new.png"`
	IsActive    *bool   `json:"is_active,omitempty" example:"true"`
}

type RewardResponse struct {
	ID          uint64     `json:"id" example:"1"`
	Name        string     `json:"name" example:"Voucher Belanja Rp 50.000"`
	Description *string    `json:"description,omitempty" example:"Voucher belanja dapat digunakan di minimarket rekanan"`
	PointCost   uint       `json:"point_cost" example:"500"`
	Stock       int        `json:"stock" example:"20"`
	Image       *string    `json:"image,omitempty" example:"rewards/voucher50k.png"`
	IsActive    bool       `json:"is_active" example:"true"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
