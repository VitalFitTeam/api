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
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	otp "github.com/vitalfit/api/pkg/otp"
)

type AuthHandlersInterface interface {
	AuthRoutes(rg *gin.RouterGroup)
	UserRoutes(rg *gin.RouterGroup)
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
		return
	}
	hash := sha256.Sum256([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])
	if err := h.services.AuthServices.RegisterUserClient(ctx, user, hashedKey); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.BadRequestResponse(c, err)
		case shared_errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	status, err := h.registerEmail(ctx, user, key, c)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
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
// @Router			/user/register-staff [post]
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
		case shared_errors.ErrNotFound:
			h.services.LogErrors.BadRequestResponse(c, err)
		case shared_errors.ErrConflict:
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

// @Summary		Activate user account
// @Description	Activates a user's account using the invitation code/token.
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			payload	body	authdomain.CodePayload	true	"Activation Code"
// @Success		204		"User successfully activated. No content returned."
// @Failure		400		{object}	map[string]interface{}	"Bad request (e.g., invalid JSON payload)"
// @Failure		404		{object}	map[string]interface{}	"Code is invalid or expired (handled by the service layer returning ErrNotFound)"
// @Failure		500		{object}	map[string]interface{}	"Internal server error (e.g., database connection issue)"
// @Router			/auth/activate [put]
func (h *AuthHandlers) ActivateUserHandler(c *gin.Context) {
	var payload authdomain.CodePayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	if err := h.services.AuthServices.Activate(ctx, payload.Code); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

// @Summary		Logs in a user and issues a JWT token
// @Description	Authenticates the user with email and password, returning an access token upon success.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		authdomain.CreateUserTokenPayload	true		"User login credentials (email and password)"
// @Success		200			{object}	map[string]string					"token"		"Successfully generated JWT access token"
// @Failure		400			{object}	map[string]string					"error":	"Invalid request body"
// @Failure		401			{object}	map[string]string					"error":	"Unauthorized"						"Invalid credentials (password mismatch)"
// @Failure		404			{object}	map[string]string					"error":	"not found"							"User with the given email not found"
// @Failure		500			{object}	map[string]string					"error":	"the server encountered a problem"	"Internal server error during token generation or hashing"
// @Router			/auth/login [post]
func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	var payload authdomain.CreateUserTokenPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user, err := h.services.UserServices.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	match, err := user.PasswordHash.Matches(payload.Password)
	if err != nil || !match {
		h.services.LogErrors.UnauthorizedErrorResponse(c, err)
		return
	}

	token, err := h.services.AuthServices.GenerateToken(user)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

// @Summary		Get current user profile
// @Description	Retrieves the profile of the user authenticated via the JWT token in the request header.
// @Tags			User
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"user"		"Current authenticated user profile"
// @Failure		401	{object}	map[string]string		"error":	"Unauthorized"	"Missing or invalid JWT token"
// @Router			/user/whoami [get]
func (h *AuthHandlers) whoami(c *gin.Context) {
	user := h.services.UserServices.GetUserFromContext(c)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *AuthHandlers) registerEmail(ctx context.Context, user *authdomain.Users, key string, c *gin.Context) (int, error) {
	status, err := h.services.AuthServices.MailSender(ctx, user, key)
	if err != nil {
		h.services.Logger.Errorw("error sending welcome email", "error", err)

		if err := h.services.AuthServices.Delete(ctx, user.UserID); err != nil {
			h.services.Logger.Errorw("error deleting user", "error", err)
			return http.StatusInternalServerError, err
		}
		h.services.LogErrors.InternalServerError(c, err)
		return http.StatusInternalServerError, err
	}

	return status, err
}
