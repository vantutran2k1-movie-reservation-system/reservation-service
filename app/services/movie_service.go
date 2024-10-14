package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/repositories"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/transaction"
	"gorm.io/gorm"
)

type MovieService interface {
	GetMovie(id uuid.UUID) (*models.Movie, *errors.ApiError)
	GetMovies(limit, offset int) ([]*models.Movie, *models.ResponseMeta, *errors.ApiError)
	CreateMovie(title string, description *string, releaseDate string, duration int, language *string, rating *float64, createdBy uuid.UUID) (*models.Movie, *errors.ApiError)
	UpdateMovie(id, updatedBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError)
	AssignGenres(id uuid.UUID, genreIDs []uuid.UUID) *errors.ApiError
}

type movieService struct {
	db                 *gorm.DB
	transactionManager transaction.TransactionManager
	movieRepo          repositories.MovieRepository
	genreRepo          repositories.GenreRepository
	movieGenreRepo     repositories.MovieGenreRepository
}

func NewMovieService(
	db *gorm.DB,
	transactionManager transaction.TransactionManager,
	movieRepo repositories.MovieRepository,
	genreRepo repositories.GenreRepository,
	movieGenreRepo repositories.MovieGenreRepository,
) MovieService {
	return &movieService{
		db:                 db,
		transactionManager: transactionManager,
		movieRepo:          movieRepo,
		genreRepo:          genreRepo,
		movieGenreRepo:     movieGenreRepo,
	}
}

func (s *movieService) GetMovie(id uuid.UUID) (*models.Movie, *errors.ApiError) {
	m, err := s.movieRepo.GetMovie(id)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if m == nil {
		return nil, errors.NotFoundError("movie not found")
	}

	return m, nil
}

func (s *movieService) GetMovies(limit, offset int) ([]*models.Movie, *models.ResponseMeta, *errors.ApiError) {
	movies, err := s.movieRepo.GetMovies(limit, offset)
	if err != nil {
		return nil, nil, errors.InternalServerError(err.Error())
	}

	count, err := s.movieRepo.GetNumbersOfMovie()
	if err != nil {
		return nil, nil, errors.InternalServerError(err.Error())
	}

	var prevUrl, nextUrl *string

	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		prevUrl = buildPaginationURL(limit, prevOffset)
	}

	if offset+limit < count {
		nextUrlOffset := offset + limit
		nextUrl = buildPaginationURL(limit, nextUrlOffset)
	}

	meta := &models.ResponseMeta{
		Limit:   limit,
		Offset:  offset,
		Total:   count,
		NextUrl: nextUrl,
		PrevUrl: prevUrl,
	}

	return movies, meta, nil
}

func (s *movieService) CreateMovie(title string, description *string, releaseDate string, duration int, language *string, rating *float64, createdBy uuid.UUID) (*models.Movie, *errors.ApiError) {
	m := models.Movie{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		ReleaseDate:     releaseDate,
		DurationMinutes: duration,
		Language:        language,
		Rating:          rating,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		CreatedBy:       createdBy,
		LastUpdatedBy:   createdBy,
	}
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.movieRepo.CreateMovie(tx, &m)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return &m, nil
}

func (s *movieService) UpdateMovie(id, updatedBy uuid.UUID, title string, description *string, releaseDate string, duration int, language *string, rating *float64) (*models.Movie, *errors.ApiError) {
	m, err := s.movieRepo.GetMovie(id)
	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}
	if m == nil {
		return nil, errors.NotFoundError("movie not found")
	}

	m.Title = title
	m.Description = description
	m.ReleaseDate = releaseDate
	m.DurationMinutes = duration
	m.Language = language
	m.Rating = rating
	m.UpdatedAt = time.Now().UTC()
	m.LastUpdatedBy = updatedBy
	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.movieRepo.UpdateMovie(tx, m)
	}); err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return m, nil
}

func (s *movieService) AssignGenres(id uuid.UUID, genreIDs []uuid.UUID) *errors.ApiError {
	m, err := s.movieRepo.GetMovie(id)
	if err != nil {
		return errors.InternalServerError(err.Error())
	}
	if m == nil {
		return errors.NotFoundError("movie not found")
	}

	allGenreIDs, err := s.genreRepo.GetGenreIDs()
	if err != nil {
		return errors.InternalServerError(err.Error())
	}

	if !allIdsInSlice(genreIDs, allGenreIDs) {
		return errors.BadRequestError("invalid genre ids")
	}

	if err := s.transactionManager.ExecuteInTransaction(s.db, func(tx *gorm.DB) error {
		return s.movieGenreRepo.UpdateGenresOfMovie(tx, id, genreIDs)
	}); err != nil {
		return errors.InternalServerError(err.Error())
	}

	return nil
}

func allIdsInSlice(first, second []uuid.UUID) bool {
	valueMap := make(map[uuid.UUID]bool)
	for _, id := range second {
		valueMap[id] = true
	}

	for _, id := range first {
		if !valueMap[id] {
			return false
		}
	}

	return true
}

func buildPaginationURL(limit, offset int) *string {
	url := fmt.Sprintf("/movies?limit=%d&offset=%d", limit, offset)
	return &url
}
