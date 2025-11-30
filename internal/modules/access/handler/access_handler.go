package accesshandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
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
