package authhandlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"github.com/vitalfit/api/internal/shared/errors"
	otp "github.com/vitalfit/api/pkg/OTP"
)

type AuthHandlersInterface interface {
	AuthRoutes(rg *gin.RouterGroup)
	RegisterUserClientHandler(c *gin.Context)
	RegisterUserStaffHandler(c *gin.Context)
}

type AuthHandlers struct {
	services appservices.Services
}

func NewAuthHandlers(services appservices.Services) *AuthHandlers {
	return &AuthHandlers{services: services}
}

// @Summary		Register New User
// @Description	Register a new user in the system with client role
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			user	body		authdomain.CreateUserClientPayload	true	"Register user data"
// @Success		201		{object}	map[string]interface{}				"message: user created"
// @Failure		400		{object}	map[string]interface{}				"bad response"
// @Failure		500		{object}	map[string]interface{}				"internal server error"
// @Router			/auth/register [post]
func (h *AuthHandlers) RegisterUserClientHandler(c *gin.Context) {
	var payload authdomain.CreateUserClientPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	birthdate, err := time.Parse(time.RFC3339, payload.BirthDate)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := &authdomain.Users{
		FirstName:        payload.FirstName,
		LastName:         payload.LastName,
		Email:            payload.Email,
		Phone:            payload.Phone,
		IdentityDocument: payload.IdentityDocument,
		BirthDate:        birthdate,
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	//store the user
	key, err := otp.GenerateCode(5)
	if err != nil {
		h.services.InternalServerError(c, err)
	}
	hash := sha256.Sum256([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])
	if err := h.services.AuthServices.RegisterUserClient(ctx, user, hashedKey); err != nil {
		switch err {
		case errors.ErrNotFound:
			h.services.LogErrors.BadRequestResponse(c, err)
		case errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	status, err := h.registerEmail(ctx, user, key, c)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"status":  status,
		"code":    key,
	})
}

// @Summary		Register New User Staff
// @Description	Register a new user in the system with and specific role
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			user	body		authdomain.CreateUserStaffPayload	true	"Register user data"
// @Success		201		{object}	map[string]interface{}				"message: user created"
// @Failure		400		{object}	map[string]interface{}				"bad response"
// @Failure		500		{object}	map[string]interface{}				"internal server error"
// @Router			/auth/register-staff [post]
func (h *AuthHandlers) RegisterUserStaffHandler(c *gin.Context) {
	var payload authdomain.CreateUserStaffPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	birthdate, err := time.Parse(time.RFC3339, payload.BirthDate)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := &authdomain.Users{
		FirstName:        payload.FirstName,
		LastName:         payload.LastName,
		Email:            payload.Email,
		Phone:            payload.Phone,
		IdentityDocument: payload.IdentityDocument,
		BirthDate:        birthdate,
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	//store the user
	key, err := otp.GenerateCode(5)
	if err != nil {
		h.services.InternalServerError(c, err)
	}
	hash := sha256.Sum256([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])
	if err := h.services.AuthServices.RegisterUserStaff(ctx, user, hashedKey, payload.RoleName); err != nil {
		switch err {
		case errors.ErrNotFound:
			h.services.LogErrors.BadRequestResponse(c, err)
		case errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	//send main
	status, err := h.registerEmail(ctx, user, key, c)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"status":  status,
		"code":    key,
	})
}

func (h *AuthHandlers) registerEmail(ctx context.Context, user *authdomain.Users, key string, c *gin.Context) (int, error) {
	status, err := h.services.AuthServices.MailSender(ctx, user, key)
	if err != nil {
		h.services.Logger.Errorw("error sending welcome email", "error", err)

		if err := h.services.AuthServices.Delete(ctx, user.UserID); err != nil {
			h.services.Logger.Errorw("error deleting user", "error", err)
		}
		h.services.LogErrors.InternalServerError(c, err)
		return http.StatusInternalServerError, err
	}

	return status, err

}
