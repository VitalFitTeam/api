package bookingdomain

import (
	"context"
	"time"

	"github.com/google/uuid"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
)

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, classID uuid.UUID) (uuid.UUID, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID) error
	GetClientSchedule(ctx context.Context, branchID uuid.UUID, userID uuid.UUID) ([]scheduledomain.Class, error)
	GetClientBookings(ctx context.Context, userID uuid.UUID) ([]BookingWithClassInfo, error)
	GetClientActualBook(ctx context.Context, userID, branchID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Booking, error)
	CanAccessService(ctx context.Context, userID, branchID, serviceID uuid.UUID) (bool, error)
	CountBookingsForClass(ctx context.Context, classID uuid.UUID) (int64, error)
	GetBookingsByClass(ctx context.Context, classID uuid.UUID) ([]*BookingWithUserInfo, error)
}
