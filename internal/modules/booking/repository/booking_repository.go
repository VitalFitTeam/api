package bookingrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	bookingdomain "github.com/vitalfit/api/internal/modules/booking/domain"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type BookingStore struct {
	db *gorm.DB
}

//
// ------------------------------------------------------------
// CreateBooking
// ------------------------------------------------------------
//

func NewBookingStore(db *gorm.DB) *BookingStore {
	return &BookingStore{db: db}
}

func (r *BookingStore) CreateBooking(ctx context.Context, booking *bookingdomain.Booking) (uuid.UUID, error) {
	if err := r.db.WithContext(ctx).Create(booking).Error; err != nil {
		return uuid.Nil, err
	}
	return booking.BookingID, nil
}

//
// ------------------------------------------------------------
// GetBookingByID
// ------------------------------------------------------------
//

func (r *BookingStore) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*bookingdomain.Booking, error) {
	var booking bookingdomain.Booking
	err := r.db.WithContext(ctx).
		Joins("Class").
		First(&booking, "booking_id = ?", bookingID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &booking, nil
}

//
// ------------------------------------------------------------
// CancelBooking
// ------------------------------------------------------------
//

func (s *BookingStore) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	return db.WithTX(s.db, func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).
			Model(&bookingdomain.Booking{}).
			Where("booking_id = ?", bookingID).
			Update("status", bookingdomain.BookingStatusCancelledByUser)

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
// CancelBookingAndUpdateBalance
// ------------------------------------------------------------
//

func (s *BookingStore) CancelBookingAndUpdateBalance(ctx context.Context, booking *bookingdomain.Booking, shouldRefundBalance bool) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if shouldRefundBalance && booking.Class.ClassID != uuid.Nil {
			result := tx.Model(&productsdomain.ClientServiceBalance{}).
				Where("user_id = ? AND service_id = ?", booking.UserID, booking.Class.ServiceID).
				Update("balance", gorm.Expr("balance + 1"))

			if result.Error != nil {
				return result.Error
			}
		}

		if err := tx.Model(&bookingdomain.Booking{}).
			Where("booking_id = ?", booking.BookingID).
			Update("status", bookingdomain.BookingStatusCancelledByUser).Error; err != nil {
			return err
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

	err := s.db.WithContext(ctx).Model(&scheduledomain.Class{}).
		Joins("JOIN bookings ON bookings.class_id = classes.class_id").
		Where("classes.branch_id = ?", branchID).
		Where("bookings.user_id = ?", userID).
		Where("bookings.status = ?", bookingdomain.BookingStatusConfirmed).
		Where("bookings.deleted_at IS NULL").
		Preload("Service").
		Preload("Instructor.User").
		Preload("Branch").
		Find(&classes).Error

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
			Where("status = ?", bookingdomain.BookingStatusConfirmed).
			Where("deleted_at IS NULL").
			Count(&count).Error
	})

	if err != nil {
		return 0, err
	}

	return count, nil
}

//
// ------------------------------------------------------------
// GetBookingsByUser
// ------------------------------------------------------------
//

func (s *BookingStore) GetClientBookings(ctx context.Context, userID uuid.UUID) ([]bookingdomain.BookingWithClassInfo, error) {
	var results []bookingdomain.BookingWithClassInfo

	query := `
        SELECT 
            b.booking_id,
            c.class_id,
            c.starts_at,
            c.ends_at,
            s.name AS service_name,
            CONCAT(u.first_name, ' ', u.last_name) AS instructor,
            br.name AS branch_name
        FROM bookings b
        JOIN classes c ON b.class_id = c.class_id
        JOIN services s ON c.service_id = s.service_id
        JOIN instructors ins ON c.instructor_id = ins.instructor_id
        JOIN users u ON ins.user_id = u.user_id
        JOIN branch br ON c.branch_id = br.branch_id
        WHERE b.user_id = ? AND b.deleted_at IS NULL
        ORDER BY c.starts_at ASC
    `

	err := s.db.WithContext(ctx).Raw(query, userID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
