package models

import (
	"github.com/google/uuid"
	"time"
)

type UserRegistrationToken struct {
	ID         uuid.UUID `json:"-" gorm:"column:id"`
	UserID     uuid.UUID `json:"-" gorm:"column:user_id"`
	TokenValue string    `json:"token" gorm:"column:token_value"`
	IsUsed     bool      `json:"-" gorm:"column:is_used"`
	CreatedAt  time.Time `json:"-" gorm:"column:created_at"`
	ExpiresAt  time.Time `json:"-" gorm:"column:expires_at"`
}
