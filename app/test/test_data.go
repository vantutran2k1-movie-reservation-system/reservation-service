package test

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
)

func GenerateRandomUser() *models.User {
	return &models.User{
		ID:    uuid.New(),
		Email: "john.doe@example.com",
	}
}
