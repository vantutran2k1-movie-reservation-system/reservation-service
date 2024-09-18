package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/auth"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/config"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: *userService}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req payloads.CreateUserRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	u, err := c.UserService.CreateUser(req.Email, req.Password)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": map[string]any{"email": u.Email}})
}

var LoginUser = func(c *gin.Context) {
	var req payloads.LoginUserRequest
	if errs := errors.BindAndValidate(c, &req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	token, err := services.LoginUser(config.DB, config.RedisClient, req.Email, req.Password)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": map[string]any{"token": token}})
}

var LogoutUser = func(c *gin.Context) {
	tokenValue := auth.GetAuthTokenFromRequest(c.Request)
	if err := services.LogoutUser(config.DB, config.RedisClient, tokenValue); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "User logout successfully"})
}

var UpdatePassword = func(c *gin.Context) {
	var req payloads.UpdatePasswordRequest
	if errs := errors.BindAndValidate(c, &req); len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserID(c)
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdatePassword(config.DB, config.RedisClient, userID, req.NewPassword); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Password is updated successfully"})
}
