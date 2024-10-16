package models

import "github.com/google/uuid"

type State struct {
	ID        uuid.UUID `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	Code      *string   `json:"code" gorm:"column:code"`
	CountryID uuid.UUID `json:"country_id" gorm:"column:country_id"`
}
