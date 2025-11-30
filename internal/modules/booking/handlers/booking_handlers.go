package bookinghandlers

import (
	"errors"
	"net/http"

	shared_errors "github.com/vitalfit/api/internal/shared/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary		Get client schedule for a branch
// @Description	Returns the available classes for a specific client in a specific branch. If the user is a client, it returns their own schedule. If staff, the user ID must be provided in the path.
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branchId	path		string	true	"Branch UUID"
// @Param			userId		path		string	false	"Client UUID (required for non-client users)"
// @Success		200			{object}	object{data=[]scheduledomain.Class}
// @Failure		400			{object}	map[string]interface{}	"Bad Request (e.g., invalid UUID, missing userId for staff)"
// @Failure		500			{object}	map[string]interface{}	"Internal Server Error"
// @Router			/schedule/branch/{branchId}/client/{userId} [get]
func (h *BookingHandlers) GetClientScheduleHandler(c *gin.Context) {
	ctx := c.Request.Context()

	branchID, err := uuid.Parse(c.Param("branchId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, errors.New("invalid branchId format"))
		return
	}

	user := h.services.UserServices.GetUserFromContext(c)
	var targetUserID uuid.UUID

	if user.Role.Name == "client" {
		targetUserID = user.UserID
	} else {
		permission := "booking:schedule:list"
		if user.Role.Name != "super_admin" {
			ok, err := h.services.UserServices.RoleHasPermission(ctx, user.RoleID, permission)
			if err != nil {
				h.services.LogErrors.InternalServerError(c, err)
				return
			}
			if !ok {
				h.services.LogErrors.ForbiddenResponse(c)
				return
			}
		}
		userIDStr := c.Param("userId")
		if userIDStr == "" {
			h.services.LogErrors.BadRequestResponse(c, errors.New("user_id is required for non-client users"))
			return
		}
		targetUserID, err = uuid.Parse(userIDStr)
		if err != nil {
			h.services.LogErrors.BadRequestResponse(c, errors.New("invalid userId format"))
			return
		}
	}

	classes, err := h.services.BookingServices.GetClientSchedule(ctx, branchID, targetUserID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": classes})
}

// @Summary		Book class
// @Description	Books a spot for a client in a class. If the user is a client, they book for themselves. If staff, the 'user_id' in the payload is required.
// @Tags			Booking
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			classId	path		string					true	"Class UUID"
// @Param			payload	body		BookUser				true	"User ID payload (only required for staff)"
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

	user := h.services.UserServices.GetUserFromContext(c)
	var targetUserID uuid.UUID

	if user.Role.Name == "client" {
		targetUserID = user.UserID
	} else {
		permission := "booking:create"
		if user.Role.Name != "super_admin" {
			ok, err := h.services.UserServices.RoleHasPermission(ctx, user.RoleID, permission)
			if err != nil {
				h.services.LogErrors.InternalServerError(c, err)
				return
			}
			if !ok {
				h.services.LogErrors.ForbiddenResponse(c)
				return
			}
		}
		// For staff, UserID is required in payload
		if payload.UserID == uuid.Nil {
			h.services.LogErrors.BadRequestResponse(c, errors.New("user_id is required for non-client users"))
			return
		}
		targetUserID = payload.UserID
	}
	bookingID, err := h.services.BookingServices.CreateBooking(ctx, targetUserID, classID)
	if err != nil {
		switch {
		case errors.Is(err, shared_errors.ErrPayment):
			h.services.LogErrors.PaymentRequiredResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
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

	user := h.services.UserServices.GetUserFromContext(c)
	if user.Role.Name != "client" {
		permission := "booking:cancel"
		if user.Role.Name != "super_admin" {
			ok, err := h.services.UserServices.RoleHasPermission(ctx, user.RoleID, permission)
			if err != nil {
				h.services.LogErrors.InternalServerError(c, err)
				return
			}
			if !ok {
				h.services.LogErrors.ForbiddenResponse(c)
				return
			}
		}
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
// @Description	Returns all bookings for a specific client. If the user is a client, it returns their own bookings. If staff, the user ID must be provided in the path.
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			userId	path		string	false	"User UUID (required for staff)"
// @Success		200		{object}	object{data=[]BookingResponse}
// @Failure		400		{object}	map[string]interface{}	"Bad Request (e.g., invalid UUID, missing userId for staff)"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error"
// @Router			/bookings/client/{userId} [get]
func (h *BookingHandlers) GetClientBookingsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	user := h.services.UserServices.GetUserFromContext(c)
	var targetUserID uuid.UUID
	var err error

	if user.Role.Name == "client" {
		targetUserID = user.UserID
	} else {
		permission := "booking:list"
		if user.Role.Name != "super_admin" {
			ok, err := h.services.UserServices.RoleHasPermission(ctx, user.RoleID, permission)
			if err != nil {
				h.services.LogErrors.InternalServerError(c, err)
				return
			}
			if !ok {
				h.services.LogErrors.ForbiddenResponse(c)
				return
			}
		}
		userIDStr := c.Param("userId")
		if userIDStr == "" {
			h.services.LogErrors.BadRequestResponse(c, errors.New("user_id is required for non-client users"))
			return
		}
		targetUserID, err = uuid.Parse(userIDStr)
		if err != nil {
			h.services.LogErrors.BadRequestResponse(c, errors.New("invalid userId format"))
			return
		}
	}

	bookings, err := h.services.BookingServices.GetClientBookings(ctx, targetUserID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": bookings,
	})
}
