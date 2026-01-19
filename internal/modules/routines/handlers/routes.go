package handlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type RoutineHandlersInterface interface {
}

type RoutineHandlers struct {
	services appservices.Services
}

func NewRoutineHandlers(services appservices.Services) *RoutineHandlers {
	return &RoutineHandlers{services: services}
}

func (h *RoutineHandlers) RoutineRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

}
