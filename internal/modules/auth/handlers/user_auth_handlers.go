package authhandlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

	user, err := payload.CreateUser()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])
	if err := h.services.AuthServices.RegisterUserStaff(ctx, user, hashToken, payload.RoleName); err != nil {
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

	//send mail
	status, err := h.services.AuthServices.MailSenderStaff(ctx, user, plainToken, mailer.UserStaffActivate)
	if err != nil {
		h.services.Logger.Errorw("error sending activation url to email", "error", err)
		if err := h.services.AuthServices.DeleteResetToken(ctx, user.UserID); err != nil {
			h.services.Logger.Errorw("error deleting user activation token ", "error", err)
			return
		}
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
		"status":  status,
		"code":    plainToken,
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

// @Summary		Login with OAuth provider (Google, etc.)
// @Description	Authenticates a user using a Clerk session token (JWT). The backend verifies the token
// @Description	using Clerk's JWKS and extracts the user's identity securely.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			payload	body		OAuthLoginPayload		true	"OAuth session token payload"
// @Success		200		{object}	object{token=string}	"Internal JWT generated"
// @Failure		400		{object}	object{error=string}	"Invalid payload"
// @Failure		401		{object}	object{error=string}	"Invalid OAuth session token"
// @Failure		404		{object}	object{error=string}	"User not found"
// @Failure		500		{object}	object{error=string}	"Internal server error"
// @Router			/auth/oauth-login [post]
func (h *AuthHandlers) OAuthLoginHandler(c *gin.Context) {
	var payload OAuthLoginPayload
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	jwksURL := h.config.Clerk.JwksURL
	if jwksURL == "" {
		h.services.LogErrors.InternalServerError(c, errors.New("CLERK_JWKS_URL not set in config"))
		return
	}

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	token, err := jwt.Parse(payload.SessionToken, jwks.Keyfunc)
	if err != nil || !token.Valid {
		h.services.LogErrors.UnauthorizedErrorResponse(c, errors.New("invalid OAuth session token"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		h.services.LogErrors.UnauthorizedErrorResponse(c, errors.New("invalid token claims"))
		return
	}

	var email string

	if emailClaim, exists := claims["email"].(string); exists && emailClaim != "" {
		email = emailClaim
	} else if primaryEmail, exists := claims["primary_email_address"].(string); exists && primaryEmail != "" {
		email = primaryEmail
	} else if emailAddresses, exists := claims["email_addresses"].([]interface{}); exists && len(emailAddresses) > 0 {
		if emailObj, ok := emailAddresses[0].(map[string]interface{}); ok {
			if emailAddr, ok := emailObj["email_address"].(string); ok {
				email = emailAddr
			}
		}
	}

	if email == "" {
		h.services.LogErrors.UnauthorizedErrorResponse(c, errors.New("email not found in token"))
		return
	}

	user, err := h.services.UserServices.GetByEmail(ctx, email)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}

	internalToken, err := h.services.AuthServices.GenerateToken(user)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": internalToken,
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

// @Summary		Get user by ID
// @Description	Retrieves the details of a specific user by their ID.
// @Tags			User
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string							true	"User ID (UUID)"
// @Success		200	{object}	object{data=GetUserResponse}	"User details response"
// @Failure		400	{object}	object{error=string}			"Bad Request: Invalid UUID format"
// @Failure		404	{object}	object{error=string}			"Not Found: User not found"
// @Failure		500	{object}	object{error=string}			"Error: Internal server error"
// @Router			/user/{id} [get]
func (h *AuthHandlers) GetUserByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user, err := h.services.UserServices.GetByID(ctx, id)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	resp := GetUserResponse{
		UserID:            user.UserID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		IdentityDocument:  user.IdentityDocument,
		BirthDate:         user.BirthDate.Format("2006-01-02"),
		Gender:            string(user.Gender),
		Phone:             user.Phone,
		ProfilePictureURL: user.ProfilePictureURL,
		RoleID:            user.RoleID,
		RoleName:          user.Role.Name,
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// @Summary		Update a staff user
// @Description	Updates the details of a non-client (staff) user.
// @Tags			User
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id		path		string					true	"User ID (UUID)"
// @Param			payload	body		UpdateUserStaffPayload	true	"User update payload"
// @Success		200		{object}	nil						"User updated successfully"
// @Failure		400		{object}	object{error=string}	"Bad Request: Invalid UUID or payload"
// @Failure		500		{object}	object{error=string}	"Error: Internal server error"
// @Router			/user/{id}/staff [put]
func (h *AuthHandlers) UpdateUserStaffHandler(c *gin.Context) {
	authenticatedUser := h.services.UserServices.GetUserFromContext(c)

	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	isOwner := authenticatedUser.UserID == id
	if authenticatedUser.Role.Name == "super_admin" {
	} else if !isOwner {
		allowed, err := h.services.UserServices.RoleHasPermission(ctx, authenticatedUser.RoleID, "users:update")
		if err != nil {
			h.services.LogErrors.InternalServerError(c, err)
			return
		}
		if !allowed {
			h.services.LogErrors.ForbiddenResponse(c)
			return
		}
	}

	var payload UpdateUserStaffPayload
	if err = c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user, err := payload.CreateUser()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user.UserID = id
	if err = h.services.UserServices.UpdateStaff(ctx, user, payload.RoleName); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Update a client user
// @Description	Updates the details of a client user.
// @Tags			User
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			id		path		string					true	"User ID (UUID)"
// @Param			payload	body		UpdateUserClientPayload	true	"User update payload"
// @Success		200		{object}	nil						"User updated successfully"
// @Failure		400		{object}	object{error=string}	"Bad Request: Invalid UUID or payload"
// @Failure		500		{object}	object{error=string}	"Error: Internal server error"
// @Router			/user/{id}/client [put]
func (h *AuthHandlers) UpdateUserClientHandler(c *gin.Context) {
	authenticatedUser := h.services.UserServices.GetUserFromContext(c)

	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	isOwner := authenticatedUser.UserID == id
	if authenticatedUser.Role.Name == "super_admin" {
	} else if !isOwner {
		allowed, err := h.services.UserServices.RoleHasPermission(ctx, authenticatedUser.RoleID, "users:update")
		if err != nil {
			h.services.LogErrors.InternalServerError(c, err)
			return
		}
		if !allowed {
			h.services.LogErrors.ForbiddenResponse(c)
			return
		}
	}

	var payload UpdateUserClientPayload
	if err = c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user, err := payload.CreateUser()
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user.UserID = id
	if err = h.services.UserServices.UpdateClient(ctx, user); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Validate Password Reset Token
// @Description	Checks if a password reset token is valid and has not expired.
// @Tags			Auth
// @Produce		json
// @Param			token	path		string					true	"Password Reset Token"
// @Success		204		{object}	nil						"Token is valid."
// @Failure		404		{object}	object{error=string}	"Token is invalid or has expired."
// @Failure		500		{object}	object{error=string}	"Internal server error."
// @Router			/auth/password/validate/{token} [get]
func (h *AuthHandlers) ValidateResetTokenHandler(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.Param("token")

	err := h.services.AuthServices.ValidateResetToken(ctx, token)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Reset token is invalid or has expired.",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "An internal server error occurred.",
			})
			return
		}
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Get user by email
// @Description	Retrieves the details of a specific user by their email address.
// @Tags			User
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		GetUserByEmailPayload			true	"User email payload"
// @Success		200		{object}	object{data=GetUserResponse}	"User details response"
// @Failure		400		{object}	object{error=string}			"Bad Request: Invalid payload"
// @Failure		403		{object}	object{error=string}			"Forbidden: Insufficient permissions"
// @Failure		404		{object}	object{error=string}			"Not Found: User not found"
// @Failure		500		{object}	object{error=string}			"Error: Internal server error"
// @Router			/user/by-email [post]
func (h *AuthHandlers) GetUserByEmailHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var payload GetUserByEmailPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	email := payload.Email
	user, err := h.services.UserServices.GetByEmail(ctx, email)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
			return
		default:
			h.services.LogErrors.InternalServerError(c, err)
			return
		}
	}
	resp := GetUserResponse{
		UserID:            user.UserID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		IdentityDocument:  user.IdentityDocument,
		RoleID:            user.RoleID,
		RoleName:          user.Role.Name,
		BirthDate:         user.BirthDate.Format("2006-01-02"),
		Gender:            string(user.Gender),
		Phone:             user.Phone,
		ProfilePictureURL: user.ProfilePictureURL,
	}
	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// @Summary		Delete a user
// @Description	Deletes a user from the system. Requires appropriate permissions.
// @Tags			User
// @Security		ApiKeyAuth
// @Produce		json
// @Param			id	path		string					true	"User ID (UUID)"
// @Success		204	{object}	nil						"User deleted successfully"
// @Failure		400	{object}	object{error=string}	"Bad Request: Invalid UUID format"
// @Failure		403	{object}	object{error=string}	"Forbidden: Insufficient permissions"
// @Failure		404	{object}	object{error=string}	"Not Found: User not found"
// @Failure		500	{object}	object{error=string}	"Error: Internal server error"
// @Router			/user/{id} [delete]
func (h *AuthHandlers) DeleteUserHandler(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	_, err = h.services.UserServices.GetByID(ctx, id)
	if err != nil {
		switch err {
		case shared_errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	if err := h.services.UserServices.Delete(ctx, id); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary		Generate QR JWT Token
// @Description	Generates a short-lived JWT token intended for use in a QR code for authentication purposes (e.g., gym access).
// @Tags			User
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	object{token=string}	"Successfully generated QR token"
// @Failure		401	{object}	object{error=string}	"Unauthorized"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/user/qr-token [get]
func (h *AuthHandlers) GenerateQrJwtTokenHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := h.services.UserServices.GetUserFromContext(c)
	token, err := h.services.AuthServices.GenerateQrJwtToken(ctx, user)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
