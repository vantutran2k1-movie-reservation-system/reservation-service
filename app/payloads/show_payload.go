package payloads

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"time"
)

type CreateShowRequest struct {
	MovieId   uuid.UUID            `json:"movie_id" binding:"required"`
	TheaterId uuid.UUID            `json:"theater_id" binding:"required"`
	StartTime time.Time            `json:"start_time" binding:"required"`
	EndTime   time.Time            `json:"end_time" binding:"required"`
	Status    constants.ShowStatus `json:"status" binding:"required,oneof=ACTIVE CANCELLED COMPLETED EXPIRED SCHEDULED ON-HOLD"`
}
