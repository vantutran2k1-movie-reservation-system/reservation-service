package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
)

type UserProfileController struct {
	UserProfileService services.UserProfileService
}

func NewUserProfileController(userProfileService *services.UserProfileService) *UserProfileController {
	return &UserProfileController{UserProfileService: *userProfileService}
}

func (c *UserProfileController) GetProfileByUserID(ctx *gin.Context) {
	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.GetProfileByUserID(userID)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": buildUserProfileResponseData(p)})
}

func (c *UserProfileController) CreateUserProfile(ctx *gin.Context) {
	var req payloads.CreateUserProfileRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.CreateUserProfile(userID, req.FirstName, req.LastName, req.PhoneNumber, req.DateOfBirth)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": buildUserProfileResponseData(p)})
}

func (c *UserProfileController) UpdateUserProfile(ctx *gin.Context) {
	var req payloads.CreateUserProfileRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	p, err := c.UserProfileService.UpdateUserProfile(userID, req.FirstName, req.LastName, req.PhoneNumber, req.DateOfBirth)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": buildUserProfileResponseData(p)})
}

func (c *UserProfileController) UpdateProfilePicture(ctx *gin.Context) {
	userID, err := middlewares.GetUserID(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	files, err := middlewares.GetUploadedFiles(ctx, constants.PROFILE_PICTURE_REQUEST_FORM_KEY)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	if err := c.UserProfileService.UpdateProfilePicture(userID, files[0]); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Profile picture is updated successfully"})
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
