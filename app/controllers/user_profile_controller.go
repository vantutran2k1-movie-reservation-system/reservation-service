package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

var CreateUserProfile = func(c *gin.Context) {
	var req payloads.CreateUserProfileRequest
	if errs := errors.BindAndValidate(c, &req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserID(c)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := services.CreateUserProfile(config.DB, userID, req.FirstName, req.LastName, req.PhoneNumber, req.DateOfBirth)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": buildUserProfileResponseData(p)})
}

var UpdateUserProfile = func(c *gin.Context) {
	var req payloads.UpdateUserProfileRequest
	if errs := errors.BindAndValidate(c, &req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserID(c)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := services.UpdateUserProfile(config.DB, userID, req.FirstName, req.LastName, req.PhoneNumber, req.DateOfBirth)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": buildUserProfileResponseData(p)})
}

func buildUserProfileResponseData(p *models.UserProfile) map[string]any {
	data := map[string]any{
		"first_name": p.FirstName,
		"last_name":  p.LastName,
	}
	if p.PhoneNumber != nil {
		data["phone_number"] = p.PhoneNumber
	}
	if p.DateOfBirth != nil {
		data["date_of_birth"] = p.DateOfBirth
	}

	return data
}
