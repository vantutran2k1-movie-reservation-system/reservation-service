package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"net/http"
	"strconv"
)

type MovieController struct {
	MovieService services.MovieService
}

func NewMovieController(movieService *services.MovieService) *MovieController {
	return &MovieController{MovieService: *movieService}
}

func (c *MovieController) GetMovie(ctx *gin.Context) {
	id, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	includeGenre := ctx.Query(constants.IncludeGenres) == "true"
	m, err := c.MovieService.GetMovie(id, includeGenre)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(m)})
}

func (c *MovieController) GetMovies(ctx *gin.Context) {
	limitParam := ctx.DefaultQuery("limit", "10")
	limit, e := strconv.Atoi(limitParam)
	if e != nil || limit <= 0 {
		limit = 10
	}

	offsetParam := ctx.DefaultQuery("offset", "0")
	offset, e := strconv.Atoi(offsetParam)
	if e != nil || offset < 0 {
		offset = 0
	}

	movies, meta, err := c.MovieService.GetMovies(limit, offset)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.SliceToMaps(movies), "meta": utils.StructToMap(meta)})
}

func (c *MovieController) CreateMovie(ctx *gin.Context) {
	var req payloads.CreateMovieRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	m, err := c.MovieService.CreateMovie(req, s.UserID)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(m)})
}

func (c *MovieController) UpdateMovie(ctx *gin.Context) {
	var req payloads.UpdateMovieRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	s, err := middlewares.GetUserSession(ctx)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	id, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	m, err := c.MovieService.UpdateMovie(id, s.UserID, req)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": utils.StructToMap(m)})
}

func (c *MovieController) UpdateMovieGenres(ctx *gin.Context) {
	id, e := uuid.Parse(ctx.Param("id"))
	if e != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var req payloads.UpdateMovieGenresRequest
	if errs := errors.BindAndValidate(ctx, &req); len(errs) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	if err := c.MovieService.AssignGenres(id, req.GenreIDs); err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "genres of movie are updated successfully"})
}
