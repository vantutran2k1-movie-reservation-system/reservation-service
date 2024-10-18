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

type StateController struct {
	StateService services.StateService
}

func NewStateController(stateService *services.StateService) *StateController {
	return &StateController{
		StateService: *stateService,
	}
}

func (c *StateController) GetStatesByCountry(ctx *gin.Context) {
	countryID, e := uuid.Parse(ctx.Param("countryId"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid country id"})
		return
	}

	states, err := c.StateService.GetStatesByCountry(countryID)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(states)})
}

func (c *StateController) CreateState(ctx *gin.Context) {
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

	state, err := c.StateService.CreateState(countryID, req.Name, req.Code)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(state)})
}
