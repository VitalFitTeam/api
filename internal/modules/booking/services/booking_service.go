package bookingservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	bookingdomain "github.com/vitalfit/api/internal/modules/booking/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
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

	if class.StartsAt.Before(time.Now()) {
		return uuid.Nil, shared_errors.ErrPastClass
	}

	if class.MaxCapacity > 0 {
		count, err := s.store.Booking.CountBookingsForClass(ctx, classID)
		if err != nil {
			return uuid.Nil, err
		}
		if count >= int64(class.MaxCapacity) {
			return uuid.Nil, shared_errors.ErrFullClass
		}
	}

	isMember, err := s.store.Membership.ClientHasActiveMembership(ctx, userID)
	if err != nil {
		return uuid.Nil, err
	}

	if isMember {
		branchService, err := s.store.Products.GetBranchServiceByID(ctx, class.BranchID, class.ServiceID)
		if err != nil {
			return uuid.Nil, err
		}

		if branchService.PriceForMember == 0 {
			booking := &bookingdomain.Booking{
				UserID:  userID,
				ClassID: classID,
				Status:  bookingdomain.BookingStatusConfirmed,
			}
			return s.store.Booking.CreateBooking(ctx, booking)
		}
	}

	clientBalance, err := s.store.Products.GetClientBalance(ctx, userID, class.ServiceID)
	if err != nil {
		if !errors.Is(err, shared_errors.ErrNotFound) {
			return uuid.Nil, err
		}
	}

	if clientBalance != nil && clientBalance.Balance > 0 {
		err := s.store.Products.SpendClientBalance(ctx, userID, class.ServiceID)
		if err != nil {
			return uuid.Nil, err
		}
		booking := &bookingdomain.Booking{UserID: userID, ClassID: classID, Status: bookingdomain.BookingStatusConfirmed}
		return s.store.Booking.CreateBooking(ctx, booking)
	}

	return uuid.Nil, shared_errors.ErrPayment
}

//
// ------------------------------------------------------------
// CancelBooking
// ------------------------------------------------------------
//

func (s *BookingService) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	// 1. Obtener los detalles de la reserva para la lógica de negocio.
	booking, err := s.store.Booking.GetBookingByID(ctx, bookingID)
	if err != nil {
		if errors.Is(err, shared_errors.ErrNotFound) {
			return shared_errors.ErrNotFound
		}
		return err
	}
	class, err := s.store.Schedule.GetClassByID(ctx, booking.ClassID)
	if err != nil {
		return err
	}
	booking.Class = *class

	// 2. Determinar si se debe reponer el saldo del cliente.
	isMember, err := s.store.Membership.ClientHasActiveMembership(ctx, booking.UserID)
	if err != nil {
		return err
	}

	branchService, err := s.store.Products.GetBranchServiceByID(ctx, class.BranchID, class.ServiceID)
	if err != nil {
		if errors.Is(err, shared_errors.ErrNotFound) {
			return s.store.Booking.CancelBookingAndUpdateBalance(ctx, booking, false)
		}
		return err
	}

	shouldRefundBalance := !isMember || (isMember && branchService.PriceForMember > 0)

	return s.store.Booking.CancelBookingAndUpdateBalance(ctx, booking, shouldRefundBalance)
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

func (s *BookingService) GetClientBookings(ctx context.Context, userID uuid.UUID) ([]bookingdomain.BookingWithClassInfo, error) {
	return s.store.Booking.GetClientBookings(ctx, userID)
}

func (s *BookingService) GetClientActualBook(ctx context.Context, userID, branchID uuid.UUID, startsAt time.Time, endsAt time.Time) (*bookingdomain.Booking, error) {
	return s.store.Booking.GetClientActualBook(ctx, userID, branchID, startsAt, endsAt)
}

func (s *BookingService) CanAccessService(ctx context.Context, userID, branchID, serviceID uuid.UUID) (bool, error) {
	isMember, err := s.store.Membership.ClientHasActiveMembership(ctx, userID)
	if err != nil {
		return false, err
	}

	if isMember {
		branchService, err := s.store.Products.GetBranchServiceByID(ctx, branchID, serviceID)
		if err != nil {
			if errors.Is(err, shared_errors.ErrNotFound) {
				return false, nil
			}
			return false, err
		}

		if branchService.PriceForMember == 0 {
			return true, nil
		}
	}

	clientBalance, err := s.store.Products.GetClientBalance(ctx, userID, serviceID)
	if err != nil && !errors.Is(err, shared_errors.ErrNotFound) {
		return false, err
	}

	return clientBalance != nil && clientBalance.Balance > 0, nil
}

func (s *BookingService) CountBookingsForClass(ctx context.Context, classID uuid.UUID) (int64, error) {
	return s.store.Booking.CountBookingsForClass(ctx, classID)
}

//
// ------------------------------------------------------------
// GetBookingsByClass
// ------------------------------------------------------------
//

func (s *BookingService) GetBookingsByClass(ctx context.Context, classID uuid.UUID) ([]*bookingdomain.BookingWithUserInfo, error) {
	return s.store.Booking.GetBookingsByClass(ctx, classID)
}
