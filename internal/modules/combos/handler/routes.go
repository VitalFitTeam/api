package comboshandler

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type CombosHandlerInterface interface {
	CombosRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
}

type CombosHandler struct {
	services appservices.Services
}

func NewCombosHandler(services appservices.Services) *CombosHandler {
	return &CombosHandler{services: services}
}

func (r *CombosHandler) CombosRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

}
