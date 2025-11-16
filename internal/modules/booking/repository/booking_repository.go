package bookingrepository

import (
	"context"

	"github.com/google/uuid"
	bookingdomain "github.com/vitalfit/api/internal/modules/booking/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type BookingStore struct {
	db *gorm.DB
}

func NewBookingStore(db *gorm.DB) *BookingStore {
	return &BookingStore{db: db}
}

//
// ------------------------------------------------------------
// CreateBooking
// ------------------------------------------------------------
//

func (r *BookingStore) CreateBooking(ctx context.Context, booking *bookingdomain.Booking) (uuid.UUID, error) {
	if err := r.db.WithContext(ctx).Create(booking).Error; err != nil {
		return uuid.Nil, err
	}
	return booking.BookingID, nil
}

//
// ------------------------------------------------------------
// CancelBooking
// ------------------------------------------------------------
//

func (s *BookingStore) CancelBooking(ctx context.Context, userID uuid.UUID, bookingID uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {

		result := tx.WithContext(ctx).
			Where("booking_id = ? AND user_id = ?", bookingID, userID).
			Delete(&bookingdomain.Booking{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

//
// ------------------------------------------------------------
// GetClientSchedule
// ------------------------------------------------------------
//

func (s *BookingStore) GetClientSchedule(
	ctx context.Context,
	branchID uuid.UUID,
	userID uuid.UUID,
) ([]scheduledomain.Class, error) {

	var classes []scheduledomain.Class

	err := db.WithTX(s.db, func(tx *gorm.DB) error {

		if err := tx.WithContext(ctx).
			Preload("Service").
			Preload("Instructor").
			Preload("Branch").
			Where("branch_id = ?", branchID).
			Find(&classes).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return classes, nil
}

func (s *BookingStore) CountBookingsForClass(ctx context.Context, classID uuid.UUID) (int64, error) {
	var count int64

	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		return tx.WithContext(ctx).
			Model(&bookingdomain.Booking{}).
			Where("class_id = ?", classID).
			Count(&count).Error
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}
