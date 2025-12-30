package policieshandler

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type PoliciesHandlerInterface interface {
	PolicyRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
}

type PoliciesHandler struct {
	services appservices.Services
}

func NewPoliciesHandler(services appservices.Services) *PoliciesHandler {
	return &PoliciesHandler{
		services: services,
	}
}

func (h *PoliciesHandler) PolicyRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

}
