package payloads

import (
	"github.com/google/uuid"
)

type MovieGenre struct {
	MovieId   uuid.UUID `json:"movie_id"`
	GenreId   uuid.UUID `json:"genre_id"`
	GenreName string    `json:"genre_name"`
}

type CreateMovieRequest struct {
	Title           string   `json:"title" binding:"required,min=1,max=255"`
	Description     *string  `json:"description" binding:"omitempty"`
	ReleaseDate     string   `json:"release_date" binding:"required,date"`
	DurationMinutes int      `json:"duration_minutes" binding:"required,min=1"`
	Language        *string  `json:"language" binding:"omitempty,min=1,max=50"`
	Rating          *float64 `json:"rating"  binding:"omitempty,min=0,max=5"`
	IsActive        *bool    `json:"is_active" binding:"required"`
}

type UpdateMovieRequest struct {
	Title           string   `json:"title" binding:"required,min=1,max=255"`
	Description     *string  `json:"description" binding:"omitempty"`
	ReleaseDate     string   `json:"release_date" binding:"required,date"`
	DurationMinutes int      `json:"duration_minutes" binding:"required,min=1"`
	Language        *string  `json:"language" binding:"omitempty,min=1,max=50"`
	Rating          *float64 `json:"rating"  binding:"omitempty,min=0,max=5"`
	IsActive        *bool    `json:"is_active" binding:"required"`
}

type UpdateMovieGenresRequest struct {
	GenreIDs []uuid.UUID `json:"genre_ids" binding:"required"`
}
