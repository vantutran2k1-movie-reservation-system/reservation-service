package models

type UserLocation struct {
	Latitude  float64 `json:"lat" gorm:"column:latitude"`
	Longitude float64 `json:"lon" gorm:"column:longitude"`
}
