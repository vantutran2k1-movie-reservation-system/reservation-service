package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID    `json:"id" gorm:"column:id"`
	Email        string       `json:"email" gorm:"column:email"`
	PasswordHash string       `json:"-" gorm:"column:password_hash"`
	CreatedAt    time.Time    `json:"-" gorm:"column:created_at"`
	UpdatedAt    time.Time    `json:"-" gorm:"column:updated_at"`
	Profile      *UserProfile `json:"profile,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
