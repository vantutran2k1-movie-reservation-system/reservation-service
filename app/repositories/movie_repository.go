package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"gorm.io/gorm"
	"time"
)

type MovieRepository interface {
	GetMovie(filter filters.MovieFilter, includeGenres bool) (*models.Movie, error)
	GetMovies(filter filters.MovieFilter) ([]*models.Movie, error)
	GetMoviesWithGenres(filter filters.MovieFilter) ([]*models.Movie, error)
	GetNumbersOfMovie(filter filters.MovieFilter) (int, error)
	CreateMovie(tx *gorm.DB, movie *models.Movie) error
	UpdateMovie(tx *gorm.DB, movie *models.Movie) error
	DeleteMovie(tx *gorm.DB, movie *models.Movie, deletedBy uuid.UUID) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) GetMovie(filter filters.MovieFilter, includeGenres bool) (*models.Movie, error) {
	query := filter.GetFilterQuery(r.db)
	if includeGenres {
		query = query.Preload("Genres")
	}

	var m models.Movie
	if err := query.First(&m).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &m, nil
}

func (r *movieRepository) GetMovies(filter filters.MovieFilter) ([]*models.Movie, error) {
	var movies []*models.Movie
	if err := filter.GetFilterQuery(r.db).Find(&movies).Error; err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *movieRepository) GetMoviesWithGenres(filter filters.MovieFilter) ([]*models.Movie, error) {
	var movies []*models.Movie
	if err := filter.GetFilterQuery(r.db).Find(&movies).Error; err != nil {
		return nil, err
	}

	if len(movies) == 0 {
		return movies, nil
	}

	movieIds := make([]uuid.UUID, 0, len(movies))
	for _, movie := range movies {
		movieIds = append(movieIds, movie.ID)
	}

	var movieGenres []*payloads.MovieGenre
	if err := r.db.Table("genres").
		Select("genres.id AS genre_id, genres.name AS genre_name, movie_genres.movie_id").
		Joins("JOIN movie_genres ON movie_genres.genre_id = genres.id").
		Where("movie_genres.movie_id IN ?", movieIds).
		Scan(&movieGenres).Error; err != nil {
		return nil, err
	}

	movieGenresMap := make(map[uuid.UUID][]models.Genre)
	for _, movieGenre := range movieGenres {
		genres, exists := movieGenresMap[movieGenre.MovieId]
		if !exists {
			movieGenresMap[movieGenre.MovieId] = []models.Genre{}
		}

		movieGenresMap[movieGenre.MovieId] = append(genres, models.Genre{ID: movieGenre.GenreId, Name: movieGenre.GenreName})
	}

	for _, movie := range movies {
		genres, exists := movieGenresMap[movie.ID]
		if !exists {
			movie.Genres = make([]models.Genre, 0)
		} else {
			movie.Genres = genres
		}
	}

	return movies, nil
}

func (r *movieRepository) GetNumbersOfMovie(filter filters.MovieFilter) (int, error) {
	var count int64
	if err := filter.GetFilterQuery(r.db).Model(&models.Movie{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *movieRepository) CreateMovie(tx *gorm.DB, movie *models.Movie) error {
	return tx.Create(movie).Error
}

func (r *movieRepository) UpdateMovie(tx *gorm.DB, movie *models.Movie) error {
	return tx.Save(movie).Error
}

func (r *movieRepository) DeleteMovie(tx *gorm.DB, movie *models.Movie, deletedBy uuid.UUID) error {
	return tx.Model(movie).Updates(map[string]any{"is_deleted": true, "updated_at": time.Now().UTC(), "last_updated_by": deletedBy}).Error
}
