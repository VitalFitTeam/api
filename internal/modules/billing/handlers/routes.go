package billinghandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type BillingHandlersInterface interface {
	BillingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
}

type BillingHandlers struct {
	services appservices.Services
}

func NewBillingHandlers(services appservices.Services) *BillingHandlers {
	return &BillingHandlers{
		services: services,
	}
}

func (r *BillingHandlers) BillingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

}
