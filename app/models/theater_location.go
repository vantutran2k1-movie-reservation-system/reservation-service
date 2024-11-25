package models

import "github.com/google/uuid"

type TheaterLocation struct {
	ID         uuid.UUID  `json:"id" gorm:"column:id"`
	TheaterID  *uuid.UUID `json:"theater_id,omitempty" gorm:"column:theater_id"`
	CityID     uuid.UUID  `json:"city_id" gorm:"column:city_id"`
	Address    string     `json:"address" gorm:"column:address"`
	PostalCode string     `json:"postal_code" gorm:"column:postal_code"`
	Latitude   float64    `json:"latitude" gorm:"column:latitude"`
	Longitude  float64    `json:"longitude" gorm:"column:longitude"`
}
