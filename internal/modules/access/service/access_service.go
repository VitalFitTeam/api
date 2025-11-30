package accessservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	accessdomain "github.com/vitalfit/api/internal/modules/access/domain"
	bookingservice "github.com/vitalfit/api/internal/modules/booking/services"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/internal/store"
)

type AccessService struct {
	store          store.Storage
	bookingservice bookingservice.BookingService
}

func NewAccessServices(store store.Storage, bookingservice bookingservice.BookingService) *AccessService {
	return &AccessService{
		store:          store,
		bookingservice: bookingservice,
	}
}

func (s *AccessService) ProcessCheckIn(ctx context.Context, userID, branchID uuid.UUID) (*accessdomain.CheckInResponse, error) {
	now := time.Now()

	bookingStartTime := now.Add(-15 * time.Minute)
	bookingEndTime := now.Add(15 * time.Minute)

	booking, err := s.store.Booking.GetClientActualBook(ctx, userID, branchID, bookingStartTime, bookingEndTime)
	if err != nil && !errors.Is(err, shared_errors.ErrNotFound) {
		return nil, err
	}
	fmt.Println(booking)
	if booking != nil {
		class, err := s.store.Schedule.GetClassByID(ctx, booking.ClassID)
		if err != nil {
			return nil, err
		}

		attendanceLog := &accessdomain.AttendanceLog{
			UserID:    userID,
			ClassID:   &booking.ClassID,
			ServiceID: class.ServiceID,
		}
		if err := s.store.Access.LogAttendance(ctx, attendanceLog); err != nil {
			return nil, err
		}

		return &accessdomain.CheckInResponse{
			Message:     "Welcome",
			AccessType:  "Class Reservation",
			ServiceName: class.Service.Name,
			CheckInTime: now,
		}, nil
	}

	walkInStartTime := now.Add(-10 * time.Minute)
	walkInEndTime := now.Add(10 * time.Minute)

	availableClasses, err := s.store.Schedule.GetAvailableClassesForBranch(ctx, branchID, walkInStartTime, walkInEndTime)
	if err != nil {
		return nil, err
	}

	for _, class := range availableClasses {
		bookingID, err := s.bookingservice.CreateBooking(ctx, userID, class.ClassID)
		if err == nil && bookingID != uuid.Nil {

			attendanceLog := &accessdomain.AttendanceLog{
				UserID:    userID,
				ClassID:   &class.ClassID,
				ServiceID: class.ServiceID, // Asumiendo que class.ServiceID es uuid.UUID
			}
			if err := s.store.Access.LogAttendance(ctx, attendanceLog); err != nil {
				return nil, err
			}

			return &accessdomain.CheckInResponse{
				Message:     "Welcome",
				AccessType:  "Class Walk-in",
				ServiceName: class.Service.Name,
				CheckInTime: now,
			}, nil
		}
	}

	openGymService, err := s.store.Products.GetBranchServiceByName(ctx, branchID, "Open Gym")
	if err != nil {
		if errors.Is(err, shared_errors.ErrNotFound) {
			return nil, errors.New("access denied: No classes available and no Open Gym access at this branch")
		}
		return nil, err
	}

	isMember, err := s.store.Membership.ClientHasActiveMembership(ctx, userID)
	if err != nil {
		return nil, err
	}

	if isMember && openGymService.PriceForMember == 0 {
		attendanceLog := &accessdomain.AttendanceLog{
			UserID:    userID,
			ServiceID: openGymService.ServiceID,
		}
		if err := s.store.Access.LogAttendance(ctx, attendanceLog); err != nil {
			return nil, err
		}

		return &accessdomain.CheckInResponse{
			Message:     "Welcome to the Gym",
			AccessType:  "Open Gym",
			ServiceName: "Open Gym",
			CheckInTime: now,
		}, nil
	}

	clientBalance, err := s.store.Products.GetClientBalance(ctx, userID, openGymService.ServiceID)
	if err != nil && !errors.Is(err, shared_errors.ErrNotFound) {
		return nil, err
	}

	if clientBalance != nil && clientBalance.Balance > 0 {
		if err := s.store.Products.SpendClientBalance(ctx, userID, openGymService.ServiceID); err != nil {
			return nil, err
		}

		attendanceLog := &accessdomain.AttendanceLog{
			UserID:    userID,
			ServiceID: openGymService.ServiceID,
		}
		if err := s.store.Access.LogAttendance(ctx, attendanceLog); err != nil {
			return nil, err
		}

		return &accessdomain.CheckInResponse{
			Message:     "Welcome to the Gym",
			AccessType:  "Open Gym",
			ServiceName: "Open Gym",
			CheckInTime: now,
		}, nil
	}

	return nil, shared_errors.ErrPayment
}
