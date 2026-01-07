package audithandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type AuditHandlersInterface interface {
	SetupRoutes(router *gin.RouterGroup, m *auth.AuthMiddleware)
}

type AuditHandlers struct {
	services appservices.Services
}

func NewAuditHandlers(services appservices.Services) *AuditHandlers {
	return &AuditHandlers{
		services: services,
	}
}

func (h *AuditHandlers) SetupRoutes(router *gin.RouterGroup, m *auth.AuthMiddleware) {

}
