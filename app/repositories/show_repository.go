package repositories

import (
	"github.com/google/uuid"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
	"time"
)

type ShowRepository interface {
	GetShow(filter filters.ShowFilter) (*models.Show, error)
	IsShowInValidTimeRange(theaterId uuid.UUID, startTime time.Time, endTime time.Time) (bool, error)
	CreateShow(tx *gorm.DB, show *models.Show) error
}

func NewShowRepository(db *gorm.DB) ShowRepository {
	return &showRepository{db: db}
}

type showRepository struct {
	db *gorm.DB
}

func (r *showRepository) GetShow(filter filters.ShowFilter) (*models.Show, error) {
	var show models.Show
	if err := filter.GetFilterQuery(r.db).First(&show).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &show, nil
}

func (r *showRepository) IsShowInValidTimeRange(theaterId uuid.UUID, startTime, endTime time.Time) (bool, error) {
	var show models.Show
	query := `
		SELECT id
		FROM shows
		WHERE TRUE
		  	AND status IN (?, ?)
			AND theater_id = ?
			AND (
			    start_time BETWEEN ? AND ?
			    OR end_time BETWEEN ? AND ?
			    OR (start_time <= ? AND end_time >= ?)
			)
	`
	if err := r.db.Raw(query, constants.Active, constants.Scheduled, theaterId, startTime, endTime, startTime, endTime, startTime, endTime).
		First(&show).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return true, nil
		}

		return false, err
	}
	
	return false, nil
}

func (r *showRepository) CreateShow(tx *gorm.DB, show *models.Show) error {
	return tx.Create(show).Error
}
