package controllers

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{UserService: *userService}
}

func (c *UserController) GetUser(ctx *gin.Context) {
	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	u, err := c.UserService.GetUser(s.UserID)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(u)})
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var req payloads.CreateUserRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	u, err := c.UserService.CreateUser(req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(u)})
}

func (c *UserController) LoginUser(ctx *gin.Context) {
	var req payloads.LoginUserRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	token, err := c.UserService.LoginUser(req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(token)})
}

func (c *UserController) LogoutUser(ctx *gin.Context) {
	t := utils.GetAuthorizationHeader(ctx.Request)
	if err := c.UserService.LogoutUser(t); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Logout user successfully"})
}

func (c *UserController) UpdateUserPassword(ctx *gin.Context) {
	var req payloads.UpdatePasswordRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	if err := c.UserService.UpdateUserPassword(s.UserID, req); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Password is updated successfully"})
}

func (c *UserController) CreatePasswordResetToken(ctx *gin.Context) {
	var req payloads.CreatePasswordResetTokenRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	t, err := c.UserService.CreatePasswordResetToken(req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(t)})
}

func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req payloads.ResetPasswordRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	if err := c.UserService.ResetUserPassword(ctx.GetHeader(constants.ResetToken), req); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Password reset successfully"})
}
