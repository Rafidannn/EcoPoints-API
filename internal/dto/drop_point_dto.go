package dto

import "time"

type CreateDropPointRequest struct {
	Name      string   `json:"name" binding:"required,min=2,max=255" example:"Drop Point Sukajadi"`
	Address   string   `json:"address" binding:"required" example:"Jl. Sukajadi No. 123, Bandung"`
	Latitude  *float64 `json:"latitude" example:"-6.89012300"`
	Longitude *float64 `json:"longitude" example:"107.61012300"`
	IsActive  *bool    `json:"is_active" example:"true"`
}

type UpdateDropPointRequest struct {
	Name      *string  `json:"name,omitempty" example:"Drop Point Sukajadi Pusat"`
	Address   *string  `json:"address,omitempty" example:"Jl. Sukajadi No. 125, Bandung"`
	Latitude  *float64 `json:"latitude,omitempty" example:"-6.89012300"`
	Longitude *float64 `json:"longitude,omitempty" example:"107.61012300"`
	IsActive  *bool    `json:"is_active,omitempty" example:"true"`
}

type DropPointResponse struct {
	ID        uint64     `json:"id" example:"1"`
	Name      string     `json:"name" example:"Drop Point Sukajadi"`
	Address   string     `json:"address" example:"Jl. Sukajadi No. 123, Bandung"`
	Latitude  *float64   `json:"latitude,omitempty" example:"-6.89012300"`
	Longitude *float64   `json:"longitude,omitempty" example:"107.61012300"`
	IsActive  bool       `json:"is_active" example:"true"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
