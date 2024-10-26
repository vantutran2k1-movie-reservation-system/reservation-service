package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"net/http"
)

type TheaterController struct {
	TheaterService services.TheaterService
}

func NewTheaterController(theaterService *services.TheaterService) *TheaterController {
	return &TheaterController{
		TheaterService: *theaterService,
	}
}

func (c *TheaterController) GetTheater(ctx *gin.Context) {
	theaterID, e := uuid.Parse(ctx.Param("theaterId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid theater id"})
		return
	}

	includeLocation := ctx.Query(constants.IncludeTheaterLocation) == "true"
	filter := payloads.GetTheaterFilter{ID: &theaterID, IncludeLocation: &includeLocation}
	theater, err := c.TheaterService.GetTheater(filter)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(theater)})
}

func (c *TheaterController) CreateTheater(ctx *gin.Context) {
	var req payloads.CreateTheaterRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	theater, err := c.TheaterService.CreateTheater(req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(theater)})
}

func (c *TheaterController) CreateTheaterLocation(ctx *gin.Context) {
	theaterID, e := uuid.Parse(ctx.Param("theaterId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid theater id"})
		return
	}

	var req payloads.CreateTheaterLocationRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	location, err := c.TheaterService.CreateTheaterLocation(theaterID, req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(location)})
}
