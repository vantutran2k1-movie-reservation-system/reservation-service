package payloads

import "github.com/google/uuid"

type GetTheaterFilter struct {
	ID              *uuid.UUID
	Name            *string
	IncludeLocation *bool
}

type CreateTheaterRequest struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}

type CreateTheaterLocationRequest struct {
	CityID     uuid.UUID `json:"city_id" binding:"required"`
	Address    string    `json:"address" binding:"required,min=2,max=255"`
	PostalCode string    `json:"postal_code" binding:"required,min=2,max=10"`
	Latitude   float64   `json:"latitude" binding:"required"`
	Longitude  float64   `json:"longitude" binding:"required"`
}
