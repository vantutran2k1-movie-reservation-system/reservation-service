package models

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
)

type Seat struct {
	Id        uuid.UUID          `json:"id" gorm:"column:id"`
	TheaterId *uuid.UUID         `json:"theater_id" gorm:"column:theater_id"`
	Row       string             `json:"row" gorm:"column:row"`
	Number    int                `json:"number" gorm:"column:number"`
	Type      constants.SeatType `json:"type" gorm:"column:type"`
}
