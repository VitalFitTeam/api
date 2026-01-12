package accesshandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/pagination"
)

// @Summary		Process User Check-In
// @Description	Processes a user's check-in attempt via a QR code JWT and branch ID. It follows a three-level priority flow: Confirmed Booking, Class Walk-in, and Open Gym access.
// @Tags			Access
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		CheckInPayload					true	"Check-In Payload with QR Token and Branch ID"
// @Success		200		{object}	accessdomain.CheckInResponse	"Access Granted"
// @Failure		400		{object}	map[string]interface{}			"Bad Request"
// @Failure		401		{object}	map[string]interface{}			"Unauthorized"
// @Failure		402		{object}	map[string]interface{}			"Payment Required"
// @Failure		403		{object}	map[string]interface{}			"Forbidden"
// @Failure		500		{object}	map[string]interface{}			"Internal Server Error"
// @Router			/access/check-in [post]
func (h *AccessHandler) CheckInHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload CheckInPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	jwtToken, err := h.services.AuthServices.ValidateToken(payload.QrJWT)
	if err != nil {
		h.services.LogErrors.UnauthorizedErrorResponse(c, err)
		return
	}

	claims, _ := jwtToken.Claims.(jwt.MapClaims)

	userID, err := uuid.Parse(claims["sub"].(string))
	if err != nil {
		h.services.LogErrors.UnauthorizedErrorResponse(c, err)
		return
	}
	branchID, err := uuid.Parse(payload.BranchID)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	resp, err := h.services.AccessServices.ProcessCheckIn(ctx, userID, branchID)
	if err != nil {
		switch {
		case errors.Is(err, shared_errors.ErrPayment):
			h.services.LogErrors.PaymentRequiredResponse(c)
		case strings.Contains(err.Error(), "access denied"):
			h.services.LogErrors.ForbiddenResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, resp)

}

// @Summary		Get Client Attendance History
// @Description	Retrieves the complete attendance history for a specific client, ordered by date with pagination
// @Tags			Access
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string	true	"Client UUID"
// @Param			start_date	query		string	false	"Start date filter (RFC3339 format)"
// @Param			end_date	query		string	false	"End date filter (RFC3339 format)"
// @Param			page		query		int		false	"Page number (default: 1)"
// @Param			limit		query		int		false	"Items per page (default: 10, max: 20)"
// @Success		200			{object}	object{data=[]AttendanceHistoryResponse}
// @Failure		400			{object}	map[string]interface{}
// @Failure		404			{object}	map[string]interface{}
// @Failure		500			{object}	map[string]interface{}
// @Router			/clients/{id}/attendance-history [get]
func (h *AccessHandler) GetClientAttendanceHistoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	// Parse pagination parameters
	fq, err := pagination.PaginatedFeedQuery{}.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	// Parse optional date filters
	var startDate, endDate *string
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate = &startDateStr
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate = &endDateStr
	}

	// Get attendance history
	attendances, total, err := h.services.AccessServices.GetClientAttendanceHistory(ctx, clientID, startDate, endDate, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	// Build response
	resp := make([]*AttendanceHistoryResponse, 0, len(attendances))
	for _, attendance := range attendances {
		item := &AttendanceHistoryResponse{
			AttendanceID: attendance.AttendanceID,
			UserID:       attendance.UserID,
			ServiceID:    attendance.ServiceID,
			CheckInTime:  attendance.CheckInTime,
			Status:       attendance.Status,
		}

		// Add user info if preloaded
		if attendance.User.UserID != uuid.Nil {
			item.UserName = attendance.User.FirstName + " " + attendance.User.LastName
		}

		// Add service info if preloaded
		if attendance.Service.ServiceID != uuid.Nil {
			item.ServiceName = attendance.Service.Name
		}

		// Add class info if preloaded
		if attendance.Class != nil && attendance.Class.ClassID != uuid.Nil {
			item.ClassName = attendance.Service.Name + " Class"
			item.ClassTime = &attendance.Class.StartsAt
		}

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, pagination.PaginatedResponseTotal[*AttendanceHistoryResponse]{
		Data:  resp,
		Count: int64(len(resp)),
		Total: total,
	})
}

// @Summary		Get Client Service Usage
// @Description	Retrieves service usage history for a specific client including branch and timestamp with pagination
// @Tags			Access
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id			path		string	true	"Client UUID"
// @Param			start_date	query		string	false	"Start date filter (RFC3339 format)"
// @Param			end_date	query		string	false	"End date filter (RFC3339 format)"
// @Param			page		query		int		false	"Page number (default: 1)"
// @Param			limit		query		int		false	"Items per page (default: 10, max: 20)"
// @Success		200			{object}	object{data=[]ServiceUsageResponse}
// @Failure		400			{object}	map[string]interface{}
// @Failure		404			{object}	map[string]interface{}
// @Failure		500			{object}	map[string]interface{}
// @Router			/clients/{id}/service-usage [get]
func (h *AccessHandler) GetClientServiceUsageHandler(c *gin.Context) {
	ctx := c.Request.Context()

	clientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	// Parse pagination parameters
	fq, err := pagination.PaginatedFeedQuery{}.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	// Parse optional date filters
	var startDate, endDate *string
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate = &startDateStr
	}
	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate = &endDateStr
	}

	// Get service usage
	serviceUsage, total, err := h.services.AccessServices.GetClientServiceUsage(ctx, clientID, startDate, endDate, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	// Build response
	resp := make([]*ServiceUsageResponse, 0, len(serviceUsage))
	for _, usage := range serviceUsage {
		item := &ServiceUsageResponse{
			AttendanceID: usage.AttendanceID,
			ServiceID:    usage.ServiceID,
			CheckInTime:  usage.CheckInTime,
			Status:       usage.Status,
		}

		// Add service info if preloaded
		if usage.Service.ServiceID != uuid.Nil {
			item.ServiceName = usage.Service.Name
		}

		// Add branch info if class is preloaded
		if usage.Class != nil && usage.Class.ClassID != uuid.Nil {
			item.BranchID = &usage.Class.BranchID
			if usage.Class.Branch.BranchID != uuid.Nil {
				item.BranchName = usage.Class.Branch.Name
			}
		}

		resp = append(resp, item)
	}

	c.JSON(http.StatusOK, pagination.PaginatedResponseTotal[*ServiceUsageResponse]{
		Data:  resp,
		Count: int64(len(resp)),
		Total: total,
	})
}
