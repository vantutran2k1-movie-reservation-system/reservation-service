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

type LocationController struct {
	LocationService services.LocationService
}

func NewLocationController(locationService *services.LocationService) *LocationController {
	return &LocationController{LocationService: *locationService}
}

func (c *LocationController) GetCountries(ctx *gin.Context) {
	countries, err := c.LocationService.GetCountries()
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"countries": utils.SliceToMaps(countries)})
}

func (c *LocationController) CreateCountry(ctx *gin.Context) {
	var req payloads.CreateCountryRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	country, err := c.LocationService.CreateCountry(req.Name, req.Code)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(country)})
}

func (c *LocationController) GetStatesByCountry(ctx *gin.Context) {
	countryID, e := uuid.Parse(ctx.Param("countryId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
		return
	}

	states, err := c.LocationService.GetStatesByCountry(countryID)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(states)})
}

func (c *LocationController) CreateState(ctx *gin.Context) {
	countryID, e := uuid.Parse(ctx.Param("countryId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
		return
	}

	var req payloads.CreateStateRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	state, err := c.LocationService.CreateState(countryID, req.Name, req.Code)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(state)})
}

func (c *LocationController) CreateCity(ctx *gin.Context) {
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

	city, err := c.LocationService.CreateCity(countryID, stateID, req.Name)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(city)})
}
