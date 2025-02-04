package models

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"time"
)

type Show struct {
	Id        uuid.UUID            `json:"id" gorm:"column:id"`
	MovieId   *uuid.UUID           `json:"movie_id" gorm:"column:movie_id"`
	TheaterId *uuid.UUID           `json:"theater_id" gorm:"column:theater_id"`
	StartTime time.Time            `json:"start_time" gorm:"column:start_time"`
	EndTime   time.Time            `json:"end_time" gorm:"column:end_time"`
	Status    constants.ShowStatus `json:"status" gorm:"column:status"`
	CreatedAt time.Time            `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time            `json:"updated_at" gorm:"column:updated_at"`
}
