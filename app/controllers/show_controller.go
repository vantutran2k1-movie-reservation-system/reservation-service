package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"net/http"
	"strconv"
)

type ShowController struct {
	ShowService services.ShowService
}

func NewShowController(showService *services.ShowService) *ShowController {
	return &ShowController{
		ShowService: *showService,
	}
}

func (c *ShowController) GetActiveShows(ctx *gin.Context) {
	limitParam := ctx.DefaultQuery(constants.Limit, "10")
	limit, e := strconv.Atoi(limitParam)
	if e != nil || limit <= 0 || limit > 10 {
		limit = 10
	}

	offsetParam := ctx.DefaultQuery(constants.Offset, "0")
	offset, e := strconv.Atoi(offsetParam)
	if e != nil || offset < 0 {
		offset = 0
	}

	shows, err := c.ShowService.GetShows(constants.Active, limit, offset)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(shows)})
}

func (c *ShowController) GetScheduledShows(ctx *gin.Context) {
	limitParam := ctx.DefaultQuery(constants.Limit, "10")
	limit, e := strconv.Atoi(limitParam)
	if e != nil || limit <= 0 || limit > 10 {
		limit = 10
	}

	offsetParam := ctx.DefaultQuery(constants.Offset, "0")
	offset, e := strconv.Atoi(offsetParam)
	if e != nil || offset < 0 {
		offset = 0
	}

	shows, err := c.ShowService.GetShows(constants.Scheduled, limit, offset)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(shows)})
}

func (c *ShowController) CreateShow(ctx *gin.Context) {
	var req payloads.CreateShowRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	show, err := c.ShowService.CreateShow(req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(show)})
}
