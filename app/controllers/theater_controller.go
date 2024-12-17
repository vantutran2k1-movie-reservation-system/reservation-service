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
	"strconv"
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
	theater, err := c.TheaterService.GetTheater(theaterID, includeLocation)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(theater)})
}

func (c *TheaterController) GetTheaters(ctx *gin.Context) {
	limitParam := ctx.DefaultQuery(constants.Limit, "10")
	limit, e := strconv.Atoi(limitParam)
	if e != nil || limit <= 0 {
		limit = 10
	}

	offsetParam := ctx.DefaultQuery(constants.Offset, "0")
	offset, e := strconv.Atoi(offsetParam)
	if e != nil || offset < 0 {
		offset = 0
	}

	includeLocation := ctx.Query(constants.IncludeTheaterLocation) == "true"

	theaters, meta, err := c.TheaterService.GetTheaters(limit, offset, includeLocation)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(theaters), "meta": utils.StructToMap(meta)})
}

func (c *TheaterController) GetNearbyTheaters(ctx *gin.Context) {
	distanceParam := ctx.DefaultQuery(constants.MaxDistance, "5")
	distance, e := strconv.ParseFloat(distanceParam, 64)
	if e != nil || distance < 0 {
		distance = 5
	}

	theaters, err := c.TheaterService.GetNearbyTheaters(distance)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(theaters)})
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

func (c *TheaterController) UpdateTheaterLocation(ctx *gin.Context) {
	theaterID, e := uuid.Parse(ctx.Param("theaterId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid theater id"})
		return
	}

	var req payloads.UpdateTheaterLocationRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	location, err := c.TheaterService.UpdateTheaterLocation(theaterID, req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(location)})
}
