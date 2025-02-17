package payloads

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
)

type GetTheaterWithLocationResult struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	LocationId uuid.UUID `json:"location_id"`
	CityId     uuid.UUID `json:"city_id"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
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

type CreateSeatPayload struct {
	Row    string             `json:"row" binding:"required,uppercase,len=1"`
	Number int                `json:"number" binding:"required,min=1,max=50"`
	Type   constants.SeatType `json:"type" binding:"required,oneof=REGULAR VIP"`
}

type UpdateTheaterLocationRequest struct {
	CityID     uuid.UUID `json:"city_id" binding:"required"`
	Address    string    `json:"address" binding:"required,min=2,max=255"`
	PostalCode string    `json:"postal_code" binding:"required,min=2,max=10"`
	Latitude   float64   `json:"latitude" binding:"required"`
	Longitude  float64   `json:"longitude" binding:"required"`
}
