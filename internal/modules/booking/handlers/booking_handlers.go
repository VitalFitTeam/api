package bookinghandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Get client schedule for a branch
// @Description	Returns the available classes for a specific client in a specific branch.
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branchId	path		string	true	"Branch UUID"
// @Param			userId		path		string	true	"Client UUID"
// @Success		200			{object}	object{data=[]scheduledomain.Class}
// @Failure		400			{object}	map[string]interface{}	"Bad Request"
// @Failure		500			{object}	map[string]interface{}	"Internal Server Error"
// @Router			/schedule/branch/{branchId}/client/{userId} [get]
func (h *BookingHandlers) GetClientScheduleHandler(c *gin.Context) {
	ctx := c.Request.Context()

	branchID, err := uuid.Parse(c.Param("branchId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid branchId format"))
		return
	}

	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid userId format"))
		return
	}

	classes, err := h.services.BookingServices.GetClientSchedule(ctx, branchID, userID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": classes})
}

// @Summary		Book class
// @Description	Books a spot for the client in a class
// @Tags			Booking
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			classId	path		string					true	"Class UUID"
// @Param			payload	body		BookUser				true	"User ID payload"
// @Success		201		{object}	BookingCreatedResponse	"Booking created"
// @Failure		400		{object}	map[string]interface{}	"Bad Request"
// @Failure		409		{object}	map[string]interface{}	"Conflict"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error"
// @Router			/schedule/{classId}/book [post]
func (h *BookingHandlers) CreateBookingHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener class_id del path
	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	var payload BookUser

	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	userID := payload.UserID

	bookingID, err := h.services.BookingServices.CreateBooking(ctx, userID, classID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, BookingCreatedResponse{
		BookingID: bookingID,
		Message:   "Booking created successfully",
	})
}

// @Summary		Cancel booking
// @Description	Cancels an existing booking made by the client
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			bookingId	path		string						true	"Booking UUID"
// @Success		200			{object}	BookingCancelledResponse	"Booking cancelled"
// @Failure		400			{object}	map[string]interface{}		"Bad Request"
// @Failure		403			{object}	map[string]interface{}		"Forbidden"
// @Failure		500			{object}	map[string]interface{}		"Internal Server Error"
// @Router			/bookings/{bookingId}/cancel [patch]
func (h *BookingHandlers) CancelBookingHandler(c *gin.Context) {
	ctx := c.Request.Context()

	bookingID, err := uuid.Parse(c.Param("bookingId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := h.services.BookingServices.CancelBooking(ctx, bookingID); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, BookingCancelledResponse{
		Message: "Booking cancelled successfully",
	})
}

// @Summary		Get client bookings
// @Description	Returns all bookings for a specific client
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			userId	path		string	true	"User UUID"
// @Success		200		{object}	object{data=[]BookingResponse}
// @Failure		400		{object}	map[string]interface{}	"Bad Request"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error"
// @Router			/bookings/client/{userId} [get]
func (h *BookingHandlers) GetClientBookingsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	bookings, err := h.services.BookingServices.GetClientBookings(ctx, userID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
	})
}
