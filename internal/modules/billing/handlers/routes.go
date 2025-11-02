package billinghandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

var validate = validator.New()

type BillingHandlersInterface interface {
	BillingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	GetPaymentMethodsHandler(c *gin.Context)
	GetPaymentMethodByIDHandler(c *gin.Context)

	CreatePaymentMethodHandler(c *gin.Context)
	UpdatePaymentMethodHandler(c *gin.Context)
	DeletePaymentMethodHandler(c *gin.Context)

	AddPaymentMethodsToBranchHandler(c *gin.Context)
	DeletePaymentMethodsFromBranchHandler(c *gin.Context)
	GetPaymentMethodsFromBranchHandler(c *gin.Context)
	UpdatePaymentMethodFromBranchHandler(c *gin.Context)
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

	branchPaymentMethodsGroup := rg.Group("/branches/:id/payment-methods")
	branchPaymentMethodsGroup.Use(m.AuthJwtTokenMiddleware())
	{
		branchPaymentMethodsGroup.GET("", m.RBACPermission("billing:list"), r.GetPaymentMethodsFromBranchHandler)
		branchPaymentMethodsGroup.POST("", m.RBACPermission("billing:create"), r.AddPaymentMethodsToBranchHandler)
		branchPaymentMethodsGroup.PUT("/:method_id", m.RBACPermission("billing:update"), r.UpdatePaymentMethodFromBranchHandler)
		branchPaymentMethodsGroup.DELETE("/:method_id", m.RBACPermission("billing:delete"), r.DeletePaymentMethodsFromBranchHandler)
	}

}
