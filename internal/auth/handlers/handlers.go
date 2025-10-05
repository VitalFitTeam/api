package authhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"github.com/vitalfit/api/internal/shared/errors"
)

type AuthHandlersInterface interface {
	RegisterUserHandler(c *gin.Context)
}

type AuthHandlers struct {
	services appservices.Services
}

type createUserPayload struct {
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Phone            string `json:"phone"`
	IdentityDocument string `json:"identity_document"`
	Password         string `json:"password" binding:"required,min=8"`
	RoleName         string `json:"role_name" binding:"required"`
}

func NewAuthHandlers(services appservices.Services) *AuthHandlers {
	return &AuthHandlers{services: services}
}

// @Summary		Register New User
// @Description	Register a new user in the system with and specific role
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			user	body		createUserPayload		true	"Register user data"
// @Success		201		{object}	map[string]interface{}	"message: user created"
// @Failure		400		{object}	map[string]interface{}	"bad response"
// @Failure		500		{object}	map[string]interface{}	"internal server error"
// @Router			/auth/register [post]
func (h *AuthHandlers) RegisterUserHandler(c *gin.Context) {
	var payload createUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}
	user := authdomain.Users{
		FirstName:        payload.FirstName,
		LastName:         payload.LastName,
		Email:            payload.Email,
		Phone:            payload.Phone,
		IdentityDocument: payload.IdentityDocument,
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}
	if err := h.services.AuthServices.RegisterUser(c.Request.Context(), user, payload.RoleName); err != nil {
		switch err {
		case errors.ErrNotFound:
			h.services.LogErrors.NotFoundResponse(c)
		case errors.ErrConflict:
			h.services.LogErrors.ConflictResponse(c, err)
		default:
			h.services.LogErrors.InternalServerError(c, err)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "user created",
	})
}
