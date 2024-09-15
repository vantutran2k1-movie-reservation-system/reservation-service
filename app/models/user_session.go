package models

import "github.com/google/uuid"

type UserSession struct {
	UserID uuid.UUID `json:"user_id"`
}
