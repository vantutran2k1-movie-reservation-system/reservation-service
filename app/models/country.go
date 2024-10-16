package models

import "github.com/google/uuid"

type Country struct {
	ID   uuid.UUID ` json:"id" gorm:"column:id"`
	Name string    `json:"name" gorm:"column:name"`
	Code string    `json:"code" gorm:"column:code"`
}
