package models

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Description     *string   `json:"description,omitempty" gorm:"default:null"`
	ReleaseDate     string    `json:"release_date"`
	DurationMinutes int       `json:"duration_minutes"`
	Language        *string   `json:"language,omitempty" gorm:"default:null"`
	Rating          *float64  `json:"rating,omitempty" gorm:"default:null"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
