package controllers

import (
	"github.com/google/uuid"
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

func (c *GenreController) GetGenre(ctx *gin.Context) {
	id, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	g, err := c.GenreService.GetGenre(id)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(g)})
}

func (c *GenreController) GetGenres(ctx *gin.Context) {
	genres, err := c.GenreService.GetGenres()
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.SliceToMaps(genres))
}

func (c *GenreController) CreateGenre(ctx *gin.Context) {
	var req payloads.CreateGenreRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	g, err := c.GenreService.CreateGenre(req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(g)})
}

func (c *GenreController) UpdateGenre(ctx *gin.Context) {
	id, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	var req payloads.UpdateGenreRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	g, err := c.GenreService.UpdateGenre(id, req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(g)})
}

func (c *GenreController) DeleteGenre(ctx *gin.Context) {
	id, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	if err := c.GenreService.DeleteGenre(id); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
