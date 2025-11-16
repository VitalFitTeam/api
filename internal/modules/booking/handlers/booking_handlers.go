package bookinghandlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ------------------------------
// GET /v1/branches/{branchId}/schedule
// ------------------------------

// @Summary		Get client schedule
// @Description	Returns the available classes for the authenticated client
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			branch_id	query		string	false	"Branch UUID"
// @Success		200	{object}	object{data=[]ScheduleClassResponse}
// @Failure		400	{object}	map[string]interface{}
// @Failure		500	{object}	map[string]interface{}
// @Router			/schedule [get]
func (h *BookingHandlers) GetClientScheduleHandler(c *gin.Context) {
	ctx := c.Request.Context()

	branchIDParam := c.Query("branch_id")
	if branchIDParam == "" {
		h.services.LogErrors.BadRequestResponse(c, errors.New("branch_id is required"))
		return
	}

	branchID, err := uuid.Parse(branchIDParam)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	userIDStr := c.GetString("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	classes, err := h.services.BookingServices.GetClientSchedule(ctx, branchID, userID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": classes})
}

// ------------------------------
// POST /v1/schedule/:classId/book
// ------------------------------

// @Summary		Book class
// @Description	Books a spot for the client in a class
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			classId	path		string	true	"Class UUID"
// @Success		201	{object}	BookingCreatedResponse
// @Failure		400	{object}	map[string]interface{}
// @Failure		409	{object}	map[string]interface{}
// @Failure		500	{object}	map[string]interface{}
// @Router			/schedule/{classId}/book [post]
func (h *BookingHandlers) CreateBookingHandler(c *gin.Context) {
	ctx := c.Request.Context()

	// Obtener class_id del path
	classID, err := uuid.Parse(c.Param("classId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	// Obtener user_id desde el body de la solicitud
	var requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	userID := requestBody.UserID

	// Crear la reserva de clase — AHORA CAPTURANDO AMBOS VALORES
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

// ------------------------------
// DELETE /v1/bookings/:bookingId
// ------------------------------

// @Summary		Cancel booking
// @Description	Cancels an existing booking made by the client
// @Tags			Booking
// @Security		ApiKeyAuth
// @Produce		json
// @Param			bookingId	path		string	true	"Booking UUID"
// @Success		200	{object}	BookingCancelledResponse
// @Failure		400	{object}	map[string]interface{}
// @Failure		403	{object}	map[string]interface{}
// @Failure		500	{object}	map[string]interface{}
// @Router			/bookings/{bookingId} [delete]
func (h *BookingHandlers) CancelBookingHandler(c *gin.Context) {
	ctx := c.Request.Context()

	bookingID, err := uuid.Parse(c.Param("bookingId"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	// Obtener user_id desde el body
	var requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	userID := requestBody.UserID

	// Cancelar la reserva
	if err := h.services.BookingServices.CancelBooking(ctx, userID, bookingID); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, BookingCancelledResponse{
		Message: "Booking cancelled successfully",
	})
}
