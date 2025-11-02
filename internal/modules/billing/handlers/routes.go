package billinghandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type BillingHandlersInterface interface {
	BillingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	GetPaymentMethodsHandler(c *gin.Context)
	GetPaymentMethodByIDHandler(c *gin.Context)

	CreatePaymentMethodHandler(c *gin.Context)
	UpdatePaymentMethodHandler(c *gin.Context)
	DeletePaymentMethodHandler(c *gin.Context)
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
	billingGroup := rg.Group("/billing")
	billingGroup.Use(m.AuthJwtTokenMiddleware())
	paymentMethodsGroup := billingGroup.Group("/payment-methods")
	{
		paymentMethodsGroup.GET("", m.RBACPermission("billing:list"), r.GetPaymentMethodsHandler)
		paymentMethodsGroup.POST("", m.RBACPermission("billing:create"), r.CreatePaymentMethodHandler)
		paymentMethodsGroup.GET("/:id", m.RBACPermission("billing:get"), r.GetPaymentMethodByIDHandler)
		paymentMethodsGroup.PUT("/:id", m.RBACPermission("billing:update"), r.UpdatePaymentMethodHandler)
		paymentMethodsGroup.DELETE("/:id", m.RBACPermission("billing:delete"), r.DeletePaymentMethodHandler)
	}
}
