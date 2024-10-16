package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"net/http"
)

type CountryController struct {
	CountryService services.CountryService
}

func NewCountryController(countryService *services.CountryService) *CountryController {
	return &CountryController{CountryService: *countryService}
}

func (c *CountryController) GetCountries(ctx *gin.Context) {
	countries, err := c.CountryService.GetCountries()
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"countries": utils.SliceToMaps(countries)})
}

func (c *CountryController) CreateCountry(ctx *gin.Context) {
	var req payloads.CreateCountryRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	country, err := c.CountryService.CreateCountry(req.Name, req.Code)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(country)})
}
