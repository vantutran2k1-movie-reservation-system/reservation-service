package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/middlewares"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/services"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
)

type MovieController struct {
	MovieService services.MovieService
}

func NewMovieController(movieService *services.MovieService) *MovieController {
	return &MovieController{MovieService: *movieService}
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

	m, err := c.MovieService.CreateMovie(req.Title, req.Description, req.ReleaseDate, req.DurationMinutes, req.Language, req.Rating, s.UserID)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie id"})
		return
	}

	m, err := c.MovieService.UpdateMovie(id, s.UserID, req.Title, req.Description, req.ReleaseDate, req.DurationMinutes, req.Language, req.Rating)
	if err != nil {
		ctx.JSON(err.StatusCode, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": utils.StructToMap(m)})
}
