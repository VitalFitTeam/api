package bookingservice

import (
	"context"
	"errors"

	"github.com/google/uuid"
	bookingdomain "github.com/vitalfit/api/internal/modules/booking/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	"github.com/vitalfit/api/internal/store"
)

type BookingService struct {
	store store.Storage
}

func NewBookingService(store store.Storage) *BookingService {
	return &BookingService{store: store}
}

//
// ------------------------------------------------------------
// CreateBooking
// ------------------------------------------------------------
//

func (s *BookingService) CreateBooking(ctx context.Context, userID uuid.UUID, classID uuid.UUID) (uuid.UUID, error) {

	class, err := s.store.Schedule.GetClassByID(ctx, classID)
	if err != nil {
		return uuid.Nil, err
	}

	if class.MaxCapacity > 0 {
		count, err := s.store.Booking.CountBookingsForClass(ctx, classID)
		if err != nil {
			return uuid.Nil, err
		}

		if count >= int64(class.MaxCapacity) {
			return uuid.Nil, errors.New("class is full")
		}
	}

	booking := &bookingdomain.Booking{
		BookingID: uuid.New(),
		UserID:    userID,
		ClassID:   classID,
		Status:    "Confirmed",
	}

	id, err := s.store.Booking.CreateBooking(ctx, booking)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

//
// ------------------------------------------------------------
// CancelBooking
// ------------------------------------------------------------
//

func (s *BookingService) CancelBooking(ctx context.Context, userID uuid.UUID, bookingID uuid.UUID) error {
	return s.store.Booking.CancelBooking(ctx, userID, bookingID)
}

//
// ------------------------------------------------------------
// GetClientSchedule
// ------------------------------------------------------------
//

func (s *BookingService) GetClientSchedule(
	ctx context.Context,
	branchID uuid.UUID,
	userID uuid.UUID,
) ([]scheduledomain.Class, error) {

	return s.store.Booking.GetClientSchedule(ctx, branchID, userID)
}

//
// ------------------------------------------------------------
// GetClientBookings
// ------------------------------------------------------------
//

// GetClientBookings devuelve todas las reservas de un usuario específico.
func (s *BookingService) GetClientBookings(ctx context.Context, userID uuid.UUID) ([]bookingdomain.BookingWithClassInfo, error) {
	return s.store.Booking.GetClientBookings(ctx, userID)
}
