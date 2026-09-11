package dto

import "time"

type CreateWasteTypeRequest struct {
	Name           string  `json:"name" binding:"required,min=2,max=255" example:"Plastik Botol Super"`
	UnitPricePerKg float64 `json:"unit_price_per_kg" binding:"required,gte=0" example:"3500.00"`
	PointsPerKg    uint    `json:"points_per_kg" binding:"required,gte=0" example:"550"`
	Description    *string `json:"description" example:"Botol plastik bersih siap daur ulang"`
	IsActive       *bool   `json:"is_active" example:"true"`
}

type UpdateWasteTypeRequest struct {
	Name           *string  `json:"name,omitempty" example:"Plastik Botol Super Edit"`
	UnitPricePerKg *float64 `json:"unit_price_per_kg,omitempty" example:"3800.00"`
	PointsPerKg    *uint    `json:"points_per_kg,omitempty" example:"600"`
	Description    *string  `json:"description,omitempty" example:"Botol plastik bersih disortir"`
	IsActive       *bool    `json:"is_active,omitempty" example:"true"`
}

type WasteTypeResponse struct {
	ID             uint64     `json:"id" example:"1"`
	Name           string     `json:"name" example:"Plastik (Botol / Ember / Kresek)"`
	UnitPricePerKg float64    `json:"unit_price_per_kg" example:"3000.00"`
	PointsPerKg    uint       `json:"points_per_kg" example:"500"`
	Description    *string    `json:"description,omitempty" example:"Botol PET bersih"`
	IsActive       bool       `json:"is_active" example:"true"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
