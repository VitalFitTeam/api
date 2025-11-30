package accesshandler

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type AccessHandlerInterface interface {
	AccessRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
}

type AccessHandler struct {
	services appservices.Services
}

func NewAccessHandler(services appservices.Services) *AccessHandler {
	return &AccessHandler{
		services: services,
	}
}

func (r *AccessHandler) AccessRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	accessGroup := rg.Group("/access")
	accessGroup.Use(m.AuthJwtTokenMiddleware())
}
