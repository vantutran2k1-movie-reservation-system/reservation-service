package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

type UserProfileController struct {
	UserProfileService services.UserProfileService
}

func NewUserProfileController(userProfileService *services.UserProfileService) *UserProfileController {
	return &UserProfileController{UserProfileService: *userProfileService}
}

func (c *UserProfileController) GetProfileByUserID(ctx *gin.Context) {
	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.GetProfileByUserID(s.UserID)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(p)})
}

func (c *UserProfileController) CreateUserProfile(ctx *gin.Context) {
	var req payloads.CreateUserProfileRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.CreateUserProfile(s.UserID, req.FirstName, req.LastName, req.PhoneNumber, req.DateOfBirth)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(p)})
}

func (c *UserProfileController) UpdateUserProfile(ctx *gin.Context) {
	var req payloads.CreateUserProfileRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.UpdateUserProfile(s.UserID, req.FirstName, req.LastName, req.PhoneNumber, req.DateOfBirth)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(p)})
}

func (c *UserProfileController) UpdateProfilePicture(ctx *gin.Context) {
	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	files, err := middlewares.GetUploadedFiles(ctx, constants.ProfilePictureRequestFormKey)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.UpdateProfilePicture(s.UserID, files[0])
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(p)})
}

func (c *UserProfileController) DeleteProfilePicture(ctx *gin.Context) {
	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	if err := c.UserProfileService.DeleteProfilePicture(s.UserID); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Profile picture is deleted successfully"})
}
