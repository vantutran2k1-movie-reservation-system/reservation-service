package models

import "github.com/google/uuid"

type Theater struct {
	ID   uuid.UUID `json:"id" gorm:"column:id"`
	Name string    `json:"name" gorm:"column:name"`
}
