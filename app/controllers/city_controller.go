package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"net/http"
)

type CityController struct {
	CityService services.CityService
}

func NewCityController(cityService *services.CityService) *CityController {
	return &CityController{
		CityService: *cityService,
	}
}

func (c *CityController) CreateCity(ctx *gin.Context) {
	countryID, e := uuid.Parse(ctx.Param("countryId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
		return
	}

	stateID, e := uuid.Parse(ctx.Param("stateId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid state id"})
		return
	}

	var req payloads.CreateCityRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	city, err := c.CityService.CreateCity(countryID, stateID, req.Name)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(city)})
}
