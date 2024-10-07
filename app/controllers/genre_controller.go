package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

type GenreController struct {
	GenreService services.GenreService
}

func NewGenreController(genreService *services.GenreService) *GenreController {
	return &GenreController{GenreService: *genreService}
}

func (c *GenreController) CreateGenre(ctx *gin.Context) {
	var req payloads.CreateGenreRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	g, err := c.GenreService.CreateGenre(req.Name)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(g)})
}
