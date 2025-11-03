package authhandlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/mailer"
	"github.com/vitalfit/api/pkg/otp"
	"github.com/vitalfit/api/pkg/pagination"
)

// @Summary		Register New User
// @Description	Register a new user in the system with client role
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			user	body		CreateUserClientPayload	true	"Register user data"
// @Success		201		{object}	map[string]interface{}	"message: user created"
// @Failure		400		{object}	map[string]interface{}	"bad response"
// @Failure		500		{object}	map[string]interface{}	"internal server error"
// @Router			/auth/register [post]
func (h *AuthHandlers) RegisterUserClientHandler(c *gin.Context) {
	var payload CreateUserClientPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user, err := payload.CreateUser()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
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
// @Param			user	body		CreateUserStaffPayload	true	"Register user data"
// @Success		201		{object}	map[string]interface{}	"message: user created"
// @Failure		400		{object}	map[string]interface{}	"bad response"
// @Failure		500		{object}	map[string]interface{}	"internal server error"
// @Router			/auth/register-staff [post]
func (h *AuthHandlers) RegisterUserStaffHandler(c *gin.Context) {
	var payload CreateUserStaffPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	if payload.RoleName == "client" {
		h.services.LogErrors.BadRequestResponse(c, shared_errors.ErrBadRequest)
		return
	}

	user, err := payload.CreateUserClientPayload.CreateUser()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
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
// @Param			payload	body	CodePayload	true	"Activation Code"
// @Success		204		"User successfully activated. No content returned."
// @Failure		400		{object}	map[string]interface{}	"Bad request (e.g., invalid JSON payload)"
// @Failure		404		{object}	map[string]interface{}	"Code is invalid or expired (handled by the service layer returning ErrNotFound)"
// @Failure		500		{object}	map[string]interface{}	"Internal server error (e.g., database connection issue)"
// @Router			/auth/activate [put]
func (h *AuthHandlers) ActivateUserHandler(c *gin.Context) {
	var payload CodePayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	if err := h.services.AuthServices.Activate(ctx, payload.Code); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

// @Summary		Activate staff user account
// @Description	Activates a staff user's account using the invitation token from the URL and sets their initial password.
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			token	path	string						true	"Activation Token"
// @Param			payload	body	UpdateStaffPasswordPayload	true	"Password Payload"
// @Success		204		"User successfully activated and password set. No content returned."
// @Failure		400		{object}	map[string]interface{}	"Bad request (e.g., invalid JSON payload, passwords do not match)"
// @Failure		404		{object}	map[string]interface{}	"Token is invalid or expired"
// @Failure		500		{object}	map[string]interface{}	"Internal server error"
// @Router			/auth/activate/{token} [put]
func (h *AuthHandlers) ActivateStaffHanlder(c *gin.Context) {
	token := c.Param("token")
	var payload UpdateStaffPasswordPayload
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	if err := h.services.AuthServices.ActivateStaff(ctx, token, payload.Password); err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

// @Summary		Logs in a user and issues a JWT token
// @Description	Authenticates the user with email and password, returning an access token upon success.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		CreateUserTokenPayload	true		"User login credentials (email and password)"
// @Success		200			{object}	map[string]string		"token"		"Successfully generated JWT access token"
// @Failure		400			{object}	map[string]string		"error":	"Invalid request body"
// @Failure		401			{object}	map[string]string		"error":	"Unauthorized"						"Invalid credentials (password mismatch)"
// @Failure		404			{object}	map[string]string		"error":	"not found"							"User with the given email not found"
// @Failure		500			{object}	map[string]string		"error":	"the server encountered a problem"	"Internal server error during token generation or hashing"
// @Router			/auth/login [post]
func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	var payload CreateUserTokenPayload
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

	// Si el contexto es 'dashboard', no permitir login de 'client'
	if payload.Context == "dashboard" && user.Role.Name == "client" {
		h.services.LogErrors.ForbiddenResponse(c)
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
func (h *AuthHandlers) WhoAmI(c *gin.Context) {
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
// @Param			email	body		ForgotPasswordPayload	true	"Estructura que contiene el correo del usuario"
// @Success		200		{object}	map[string]interface{}	"Si el correo existe, el proceso de token ha sido exitoso (por seguridad, el mensaje no confirma la existencia del correo)."
// @Failure		400		{object}	map[string]interface{}	"Bad Request - Datos de entrada inválidos (ej. formato de email incorrecto)"
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error - Error al generar el token, al acceder a la DB, o al enviar el correo."
// @Router			/auth/password/forgot [post]
func (h *AuthHandlers) ForgotPasswordHandler(c *gin.Context) {
	var payload ForgotPasswordPayload
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

// @Summary		Execute Password Reset
// @Description	Finalizes the password reset process by validating the token and updating the user's password hash.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			body	body		ResetPasswordPayload	true	"Reset token and new password data (must include password confirmation)."
// @Success		200		{object}	map[string]interface{}	"Password updated successfully."
// @Failure		400		{object}	map[string]interface{}	"Bad Request - Invalid data (expired token, token not found, or passwords do not match)."
// @Failure		500		{object}	map[string]interface{}	"Internal Server Error - Error while hashing the password or executing the DB transaction."
// @Router			/auth/password/reset [post]
func (h *AuthHandlers) ResetPasswordHandler(c *gin.Context) {
	var payload ResetPasswordPayload
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
	c.JSON(http.StatusOK, nil)
}

// function to send emails on register
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

// @Summary		Get a user list that are branch admins
// @Description	Get a user list with the branch admin role
// @Tags			User
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Success		200	{object}	object{data=[]BranchAdminResponse}	"succed response"
// @Failure		500	{object}	object{error=string}				"Error: internal server error"
// @Router			/user/branch-admins [get]
func (h *AuthHandlers) GetBranchAdminsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	fq := pagination.PaginatedFeedQuery{
		Limit:  10,
		Page:   1,
		Sort:   "asc",
		Search: "",
	}

	fq, err := fq.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	users, err := h.services.UserServices.GetBranchAdmins(ctx, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return

	}

	responseList := make([]BranchAdminResponse, 0, len(users))
	for _, user := range users {
		resp := BranchAdminResponse{
			UserID:    user.UserID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			RoleID:    user.RoleID,
			RoleName:  user.Role.Name,
		}
		responseList = append(responseList, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseList,
	})

}

// @Summary		List staff users
// @Description	Retrieves a paginated list of all users who are not clients. Supports filtering by role and searching by name or email.
// @Tags			User
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			limit	query		int							false	"Number of results per page"
// @Param			page	query		int							false	"Page number"
// @Param			sort	query		string						false	"Sort order (asc/desc)"
// @Param			search	query		string						false	"Search term for first name, last name, or email"
// @Param			role	query		string						false	"Filter by role name"
// @Success		200		{object}	object{data=[]UserResponse}	"A list of staff users"
// @Failure		400		{object}	object{error=string}		"Error: Bad Request (e.g., invalid query parameters)"
// @Failure		500		{object}	object{error=string}		"Error: Internal Server Error"
// @Router			/user/users [get]
func (h *AuthHandlers) GetUsersHandler(c *gin.Context) {
	ctx := c.Request.Context()
	fq := pagination.PaginatedFeedQuery{
		Limit:  100,
		Page:   1,
		Sort:   "asc",
		Search: "",
		Role:   "",
	}

	fq, err := fq.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	users, err := h.services.UserServices.GetUsers(ctx, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	responseList := make([]UserResponse, 0, len(users))
	for _, user := range users {
		resp := UserResponse{

			UserID:           user.UserID,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Email:            user.Email,
			RoleID:           user.RoleID,
			RoleName:         user.Role.Name,
			IdentityDocument: user.IdentityDocument,
			IsValidated:      user.IsValidated,
		}
		responseList = append(responseList, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseList,
	})
}

// @Summary		List client users
// @Description	Retrieves a paginated list of all users with the 'client' role. Supports searching by name.
// @Tags			User
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			limit	query		int							false	"Number of results per page"
// @Param			page	query		int							false	"Page number"
// @Param			sort	query		string						false	"Sort order (asc/desc)"
// @Param			search	query		string						false	"Search term for first name or last name"
// @Success		200		{object}	object{data=[]UserResponse}	"A list of client users"
// @Failure		400		{object}	object{error=string}		"Error: Bad Request (e.g., invalid query parameters)"
// @Failure		500		{object}	object{error=string}		"Error: Internal Server Error"
// @Router			/user/clients [get]
func (h *AuthHandlers) GetClientsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	fq := pagination.PaginatedFeedQuery{
		Limit:  100,
		Page:   1,
		Sort:   "asc",
		Search: "",
	}

	fq, err := fq.Parse(c.Request)
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	users, err := h.services.UserServices.GetClients(ctx, fq)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	responseList := make([]UserResponse, 0, len(users))
	for _, user := range users {
		resp := UserResponse{

			UserID:           user.UserID,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			Email:            user.Email,
			RoleID:           user.RoleID,
			RoleName:         user.Role.Name,
			IdentityDocument: user.IdentityDocument,
			IsValidated:      user.IsValidated,
		}
		responseList = append(responseList, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": responseList,
	})
}
