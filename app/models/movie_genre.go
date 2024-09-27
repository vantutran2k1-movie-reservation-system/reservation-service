package models

import "github.com/google/uuid"

type MovieGenre struct {
	MovieID uuid.UUID `json:"movie_id" gorm:"column:movie_id"`
	GenreID uuid.UUID `json:"genre_id" gorm:"column:genre_id"`
}
