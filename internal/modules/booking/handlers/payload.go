package bookinghandlers

import (
	"time"

	"github.com/google/uuid"
)

// ----------------------------------------
// RESPONSE PAYLOADS
// ----------------------------------------

type BookUser struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

type ScheduleClassResponse struct {
	ClassID           uuid.UUID `json:"class_id"`
	BranchID          uuid.UUID `json:"branch_id"`
	ServiceID         uuid.UUID `json:"service_id"`
	ServiceName       string    `json:"service_name"`
	InstructorID      uuid.UUID `json:"instructor_id"`
	InstructorName    string    `json:"instructor_name"`
	StartsAt          time.Time `json:"starts_at"`
	EndsAt            time.Time `json:"ends_at"`
	MaxCapacity       int       `json:"max_capacity"`
	RemainingCapacity int       `json:"remaining_capacity"`
	IsBooked          bool      `json:"is_booked"`
}

type BookingCreatedResponse struct {
	BookingID uuid.UUID `json:"booking_id"`
	Message   string    `json:"message"`
}

type BookingCancelledResponse struct {
	Message string `json:"message"`
}

type BookingResponse struct {
	BookingID   uuid.UUID `json:"booking_id"`
	ClassID     uuid.UUID `json:"class_id"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	ServiceName string    `json:"service_name"`
	Instructor  string    `json:"instructor"`
	BranchName  string    `json:"branch_name"`
}
