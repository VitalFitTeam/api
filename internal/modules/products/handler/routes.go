package productshandler

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type ProductsHandlerInterface interface {
	ProductsRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)

	ListServiceCategoriesHandler(c *gin.Context)
	CreateServiceHandler(c *gin.Context)
	GetServicesHandler(c *gin.Context)
	DeleteServiceHandler(c *gin.Context)
	GetServiceByIDHandler(c *gin.Context)
	UpdateServiceHandler(c *gin.Context)
}

type ProductsHandler struct {
	services appservices.Services
}

func NewProductsHandler(services appservices.Services) *ProductsHandler {
	return &ProductsHandler{services: services}
}

func (r *ProductsHandler) ProductsRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	ProductGroup := rg.Group("/services")
	{
		ProductGroup.Use(m.AuthJwtTokenMiddleware())
		ProductGroup.GET("/categories", m.RBACPermission("services:list"), r.ListServiceCategoriesHandler)
		ProductGroup.GET("/all", m.RBACPermission("services:list"), r.GetServicesHandler)

		ProductGroup.POST("", m.RBACPermission("services:create"), r.CreateServiceHandler)
		ProductGroup.GET("/:id", m.RBACPermission("services:get"), r.GetServiceByIDHandler)
		ProductGroup.DELETE("/:id", m.RBACPermission("services:delete"), r.DeleteServiceHandler)
		ProductGroup.PUT("/:id", m.RBACPermission("services:update"), r.UpdateServiceHandler)
	}
}
