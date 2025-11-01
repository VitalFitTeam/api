package productshandler

<<<<<<< HEAD
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

	AssignBranchServiceHandler(c *gin.Context)
	GetBranchServiceHandler(c *gin.Context)
	UpdateBranchServiceHandler(c *gin.Context)
	DeleteBranchServiceHandler(c *gin.Context)
	GetBranchServiceByIDHandler(c *gin.Context)
=======
import appservices "github.com/vitalfit/api/internal/app/services"

type ProductsHandlerInterface interface {
>>>>>>> dev
}

type ProductsHandler struct {
	services appservices.Services
}

func NewProductsHandler(services appservices.Services) *ProductsHandler {
	return &ProductsHandler{services: services}
}
<<<<<<< HEAD

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

	BranchServicesGroup := rg.Group("/branches/:id/services")
	{
		BranchServicesGroup.Use(m.AuthJwtTokenMiddleware())
		BranchServicesGroup.POST("", m.RBACPermission("services:create"), r.AssignBranchServiceHandler)
		BranchServicesGroup.GET("", m.RBACPermission("services:list"), r.GetBranchServiceHandler)
		BranchServicesGroup.GET("/:service_id", m.RBACPermission("services:get"), r.GetBranchServiceByIDHandler)
		BranchServicesGroup.PUT("/:service_id", m.RBACPermission("services:update"), r.UpdateBranchServiceHandler)
		BranchServicesGroup.DELETE("/:service_id", m.RBACPermission("services:delete"), r.DeleteBranchServiceHandler)
	}
}
=======
>>>>>>> dev
