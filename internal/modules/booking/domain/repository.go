package bookingdomain

import (
	"context"
	"time"

	"github.com/google/uuid"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
)

// BookingRepository define las operaciones de acceso a datos para las reservas de clases.
type BookingRepository interface {

	// CreateBooking crea una nueva reserva.
	CreateBooking(ctx context.Context, booking *Booking) (uuid.UUID, error)
	GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*Booking, error)

	// CancelBooking elimina una reserva del usuario (validación userID + bookingID).
	CancelBooking(ctx context.Context, bookingID uuid.UUID) error
	CancelBookingAndUpdateBalance(ctx context.Context, booking *Booking, shouldRefundBalance bool) error

	// GetBookingByID obtiene una reserva por su ID.

	// Contar reservas confirmadas de una clase
	CountBookingsForClass(ctx context.Context, classID uuid.UUID) (int64, error)

	// GetClientSchedule retorna las clases visibles al cliente.
	GetClientSchedule(ctx context.Context, branchID uuid.UUID, userID uuid.UUID) ([]scheduledomain.Class, error)

	// GetClientBookings retorna todas las reservas de un usuario
	GetClientBookings(ctx context.Context, userID uuid.UUID) ([]BookingWithClassInfo, error)
	GetClientActualBook(ctx context.Context, userID, branchID uuid.UUID, startsAt time.Time, endsAt time.Time) (*Booking, error)

	// GetBookingsByClass obtiene todas las reservas de una clase con información del cliente
	GetBookingsByClass(ctx context.Context, classID uuid.UUID) ([]*BookingWithUserInfo, error)
}
