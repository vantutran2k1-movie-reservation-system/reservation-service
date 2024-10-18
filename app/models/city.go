package models

import "github.com/google/uuid"

type City struct {
	ID      uuid.UUID `json:"id" gorm:"column:id"`
	Name    string    `json:"name" gorm:"column:name"`
	StateID uuid.UUID `json:"state_id" gorm:"column:state_id"`
}
