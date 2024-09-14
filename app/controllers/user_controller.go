package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
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

var LoginUser = func(c *gin.Context) {
	var req payloads.LoginUserRequest
	if errs := errors.BindAndValidate(c, &req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	token, err := services.LoginUser(config.DB, req.Email, req.Password)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": map[string]any{"token": token}})
}

var LogoutUser = func(c *gin.Context) {
	tokenValue := auth.GetAuthTokenFromRequest(c.Request)
	if err := services.LogoutUser(config.DB, tokenValue); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "User logout successfully"})
}
