package models

import (
	"time"

	"github.com/google/uuid"
)

type LoginToken struct {
	ID         uuid.UUID `json:"-" gorm:"column:id"`
	UserID     uuid.UUID `json:"-" gorm:"column:user_id"`
	TokenValue string    `json:"token" gorm:"column:token_value"`
	CreatedAt  time.Time `json:"-" gorm:"column:created_at"`
	ExpiresAt  time.Time `json:"-" gorm:"column:expires_at"`
}
