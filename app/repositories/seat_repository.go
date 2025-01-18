package repositories

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/filters"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/models"
	"gorm.io/gorm"
)

type SeatRepository interface {
	GetSeat(filter filters.SeatFilter) (*models.Seat, error)
	CreateSeat(tx *gorm.DB, seat *models.Seat) error
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}

type seatRepository struct {
	db *gorm.DB
}

func (r *seatRepository) GetSeat(filter filters.SeatFilter) (*models.Seat, error) {
	var seat models.Seat
	if err := filter.GetFilterQuery(r.db).First(&seat).Error; err != nil {
		if errors.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &seat, nil
}

func (r *seatRepository) CreateSeat(tx *gorm.DB, seat *models.Seat) error {
	return tx.Create(seat).Error
}
