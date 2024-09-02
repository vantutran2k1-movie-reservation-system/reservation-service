package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

var CreateUser = func(c *gin.Context) {
	var req payloads.CreateUserRequest
	if errs := errors.BindAndValidate(c, &req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	u, err := services.CreateUser(config.DB, req.Email, req.Password)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": map[string]any{"email": u.Email}})
}
