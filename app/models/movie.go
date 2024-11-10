package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID              uuid.UUID `json:"id" gorm:"column:id"`
	Title           string    `json:"title" gorm:"column:title"`
	Description     *string   `json:"description,omitempty" gorm:"column:description"`
	ReleaseDate     string    `json:"release_date" gorm:"column:release_date"`
	DurationMinutes int       `json:"duration_minutes" gorm:"column:duration_minutes"`
	Language        *string   `json:"language,omitempty" gorm:"column:language"`
	Rating          *float64  `json:"rating,omitempty" gorm:"column:rating"`
	IsActive        bool      `json:"is_active" gorm:"column:is_active"`
	CreatedAt       time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"-" gorm:"column:updated_at"`
	CreatedBy       uuid.UUID `json:"created_by" gorm:"column:created_by"`
	LastUpdatedBy   uuid.UUID `json:"last_updated_by" gorm:"column:last_updated_by"`
	Genres          []Genre   `json:"genres,omitempty" gorm:"many2many:movie_genres"`
}
