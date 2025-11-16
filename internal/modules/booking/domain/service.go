package bookingdomain

import (
	"context"

	"github.com/google/uuid"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
)

type BookingServiceInterface interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, classID uuid.UUID) (uuid.UUID, error)
	CancelBooking(ctx context.Context, userID uuid.UUID, bookingID uuid.UUID) error
	GetClientSchedule(ctx context.Context, branchID uuid.UUID, userID uuid.UUID) ([]scheduledomain.Class, error)
	GetClientBookings(ctx context.Context, userID uuid.UUID) ([]BookingWithClassInfo, error)
}
