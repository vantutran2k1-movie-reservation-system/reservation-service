package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	PhoneNumber       string    `json:"phone_number,omitempty" gorm:"default:null"`
	DateOfBirth       string    `json:"date_of_birth,omitempty" gorm:"default:null"`
	ProfilePictureUrl string    `json:"profile_picture_url,omitempty" gorm:"default:null"`
	Bio               string    `json:"bio,omitempty" gorm:"default:null"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
