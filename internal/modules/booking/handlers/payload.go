package bookinghandlers

import (
	"time"

	"github.com/google/uuid"
)

// ----------------------------------------
// RESPONSE PAYLOADS
// ----------------------------------------

// ScheduleClassResponse representa una clase disponible en el horario del cliente.
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

// BookingCreatedResponse representa la respuesta luego de reservar un cupo.
type BookingCreatedResponse struct {
	BookingID uuid.UUID `json:"booking_id"`
	Message   string    `json:"message"`
}

// BookingCancelledResponse representa la respuesta luego de cancelar una reserva.
type BookingCancelledResponse struct {
	Message string `json:"message"`
}

// BookingResponse representa una reserva con información detallada de la clase.
type BookingResponse struct {
	BookingID   uuid.UUID `json:"booking_id"`
	ClassID     uuid.UUID `json:"class_id"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	ServiceName string    `json:"service_name"`
	Instructor  string    `json:"instructor"`
	BranchName  string    `json:"branch_name"`
}
