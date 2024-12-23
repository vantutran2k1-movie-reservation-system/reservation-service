package models

import "github.com/google/uuid"

type Seat struct {
	Id        uuid.UUID `json:"id" gorm:"column:id"`
	TheaterId uuid.UUID `json:"theater_id" gorm:"column:theater_id"`
	Row       string    `json:"row" gorm:"column:row"`
	Number    int       `json:"number" gorm:"column:number"`
	Type      string    `json:"type" gorm:"column:type"`
}
