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
	"github.com/vitalfit/api/internal/shared/middleware/auth"
	"github.com/vitalfit/api/pkg/mailer"
	otp "github.com/vitalfit/api/pkg/otp"
)

type AuthHandlersInterface interface {
	AuthRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	UserRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
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
func (h *AuthHandlers) registerUserClientHandler(c *gin.Context) {
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
		FirstName:         payload.FirstName,
		LastName:          payload.LastName,
		Email:             payload.Email,
		Phone:             payload.Phone,
		IdentityDocument:  payload.IdentityDocument,
		BirthDate:         birthdate,
		Gender:            authdomain.GenderEnum(payload.Gender),
		ProfilePictureURL: payload.ProfilePictureURL,
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	//store the user
	key, err := otp.GenerateCode(6)
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

	status, err := h.registerEmail(ctx, user, key)
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
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			user	body		authdomain.CreateUserStaffPayload	true	"Register user data"
// @Success		201		{object}	map[string]interface{}				"message: user created"
// @Failure		400		{object}	map[string]interface{}				"bad response"
// @Failure		500		{object}	map[string]interface{}				"internal server error"
// @Router			/auth/register-staff [post]
func (h *AuthHandlers) registerUserStaffHandler(c *gin.Context) {
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
		FirstName:         payload.FirstName,
		LastName:          payload.LastName,
		Email:             payload.Email,
		Phone:             payload.Phone,
		IdentityDocument:  payload.IdentityDocument,
		BirthDate:         birthdate,
		Gender:            authdomain.GenderEnum(payload.Gender),
		ProfilePictureURL: payload.ProfilePictureURL,
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	//store the user
	key, err := otp.GenerateCode(6)
	if err != nil {
		h.services.InternalServerError(c, err)
		return
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
	status, err := h.registerEmail(ctx, user, key)
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
func (h *AuthHandlers) activateUserHandler(c *gin.Context) {
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
func (h *AuthHandlers) loginHandler(c *gin.Context) {
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

// @Summary		Solicitar token de reseteo de contraseña
// @Description	Envía un código OTP al correo electrónico proporcionado para iniciar el proceso de reseteo de contraseña.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			email	body		authdomain.ForgotPasswordPayload	true	"Estructura que contiene el correo del usuario"
// @Success		200		{object}	map[string]interface{}				"Si el correo existe, el proceso de token ha sido exitoso (por seguridad, el mensaje no confirma la existencia del correo)."
// @Failure		400		{object}	map[string]interface{}				"Bad Request - Datos de entrada inválidos (ej. formato de email incorrecto)"
// @Failure		500		{object}	map[string]interface{}				"Internal Server Error - Error al generar el token, al acceder a la DB, o al enviar el correo."
// @Router			/auth/password/forgot [post]
func (h *AuthHandlers) forgotPasswordHandler(c *gin.Context) {
	var payload authdomain.ForgotPasswordPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	//creates otp
	key, err := otp.GenerateCode(6)
	if err != nil {
		h.services.InternalServerError(c, err)
		return
	}
	hash := sha256.Sum256([]byte(key))
	hashedKey := hex.EncodeToString(hash[:])

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

	//Sends otp key to user if exists
	err = h.services.AuthServices.CreatePasswordResetToken(ctx, user.Email, hashedKey)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	//send email -> error -> rollback
	status, err := h.services.AuthServices.MailSender(ctx, user, key, mailer.UserResetPwsTemplate)
	if err != nil {
		h.services.Logger.Errorw("error sending reset token password to email", "error", err)
		if err := h.services.AuthServices.DeleteResetToken(ctx, user.UserID); err != nil {
			h.services.Logger.Errorw("error deleting user reset token password ", "error", err)
			return
		}
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(status, gin.H{
		"message": "reset key created",
		"code":    key,
	})

}

func (h *AuthHandlers) resetPasswordHandler(c *gin.Context) {
	var payload authdomain.ResetPasswordPayload
	var user authdomain.Users
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	err := h.services.AuthServices.ResetPassword(ctx, payload.Token, &user)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.BadRequestResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

}

// @Summary		Execute Password Reset
// @Description	Finalizes the password reset process by validating the token and updating the user's password hash.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		authdomain.ResetPasswordPayload	true	"Reset token and new password data (must include password confirmation)."
// @Success		200		{object}	map[string]interface{}			"Password updated successfully."
// @Failure		400		{object}	map[string]interface{}			"Bad Request - Invalid data (expired token, token not found, or passwords do not match)."
// @Failure		500		{object}	map[string]interface{}			"Internal Server Error - Error while hashing the password or executing the DB transaction."
// @Router			/auth/password/reset [post]
func (h *AuthHandlers) registerEmail(ctx context.Context, user *authdomain.Users, key string) (int, error) {
	status, err := h.services.AuthServices.MailSender(ctx, user, key, mailer.UserWelcomeTemplate)
	if err != nil {
		h.services.Logger.Errorw("error sending welcome email", "error", err)

		// Attempt to rollback user creation
		if rollbackErr := h.services.AuthServices.Delete(ctx, user.UserID); rollbackErr != nil {
			h.services.Logger.Errorw("critical: failed to rollback user creation after email failure", "original_error", err, "rollback_error", rollbackErr, "user_id", user.UserID)
		}
		return http.StatusInternalServerError, err // Return original error
	}

	return status, nil
}
