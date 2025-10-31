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
	UpdatServicetHandler(c *gin.Context)
}

type ProductsHandler struct {
	services appservices.Services
}

func NewProductsHandler(services appservices.Services) *ProductsHandler {
	return &ProductsHandler{services: services}
}

func (r *ProductsHandler) ProductsRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

}
