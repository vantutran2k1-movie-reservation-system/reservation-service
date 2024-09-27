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
	PhoneNumber       *string   `json:"phone_number,omitempty" gorm:"column:phone_number;default:null"`
	DateOfBirth       *string   `json:"date_of_birth,omitempty" gorm:"column:date_of_birth;default:null"`
	ProfilePictureUrl *string   `json:"profile_picture_url,omitempty" gorm:"column:profile_picture_url;default:null"`
	Bio               *string   `json:"bio,omitempty" gorm:"column:bio;default:null"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
}
