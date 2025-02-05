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
	GetShows(filter filters.ShowFilter) ([]*models.Show, error)
	IsShowInValidTimeRange(theaterId uuid.UUID, startTime time.Time, endTime time.Time) (bool, error)
	CreateShow(tx *gorm.DB, show *models.Show) error
	UpdateShowStatus(tx *gorm.DB, showId uuid.UUID, status constants.ShowStatus) error
	ScheduleActivateShows(tx *gorm.DB, beforeStart time.Duration) error
	ScheduleCompleteShows(tx *gorm.DB) error
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

func (r *showRepository) GetShows(filter filters.ShowFilter) ([]*models.Show, error) {
	var shows []*models.Show
	if err := filter.GetFilterQuery(r.db).Find(&shows).Error; err != nil {
		return nil, err
	}

	return shows, nil
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

func (r *showRepository) UpdateShowStatus(tx *gorm.DB, showId uuid.UUID, status constants.ShowStatus) error {
	return tx.Model(&models.Show{}).
		Where("id = ?", showId).
		Updates(map[string]any{"status": status, "updated_at": time.Now().UTC()}).
		Error
}

func (r *showRepository) ScheduleActivateShows(tx *gorm.DB, beforeStart time.Duration) error {
	currentTime := time.Now().UTC()
	maxStartTime := currentTime.Add(beforeStart)
	return tx.Model(&models.Show{}).
		Where("start_time <= ? AND status = ?", maxStartTime, constants.Scheduled).
		Updates(map[string]interface{}{"status": constants.Active, "updated_at": currentTime}).
		Error
}

func (r *showRepository) ScheduleCompleteShows(tx *gorm.DB) error {
	currentTime := time.Now().UTC()
	return tx.Model(&models.Show{}).
		Where("end_time <= ? AND status = ?", currentTime, constants.Active).
		Updates(map[string]interface{}{"status": constants.Completed, "updated_at": currentTime}).
		Error
}
