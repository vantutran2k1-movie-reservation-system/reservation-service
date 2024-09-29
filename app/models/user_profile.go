package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                uuid.UUID `json:"id" gorm:"column:id"`
	UserID            uuid.UUID `json:"user_id" gorm:"column:user_id"`
	FirstName         string    `json:"first_name" gorm:"column:first_name"`
	LastName          string    `json:"last_name" gorm:"column:last_name"`
	PhoneNumber       *string   `json:"phone_number,omitempty" gorm:"column:phone_number"`
	DateOfBirth       *string   `json:"date_of_birth,omitempty" gorm:"column:date_of_birth"`
	ProfilePictureUrl *string   `json:"profile_picture_url,omitempty" gorm:"column:profile_picture_url"`
	Bio               *string   `json:"bio,omitempty" gorm:"column:bio"`
	CreatedAt         time.Time `json:"-" gorm:"column:created_at"`
	UpdatedAt         time.Time `json:"-" gorm:"column:updated_at"`
}
