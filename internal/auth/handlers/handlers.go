package authhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
)

type AuthHandlersInterface interface {
	RegisterUserHandler(c *gin.Context)
}

type AuthHandlers struct {
	services appservices.Services
}

func NewAuthHandlers(services appservices.Services) *AuthHandlers {
	return &AuthHandlers{services: services}
}

func (h *AuthHandlers) RegisterUserHandler(c *gin.Context) {

}
