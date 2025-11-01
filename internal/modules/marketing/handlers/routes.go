package marketinghandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type MarketingHandlerInterface interface {
	MarketingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	CreateBannerHandler(c *gin.Context)
	UpdateBannerHandler(c *gin.Context)
	DeleteBannerHandler(c *gin.Context)
	GetBannerByIDHandler(c *gin.Context)
	GetBannersHandler(c *gin.Context)
}

type MarketingHandler struct {
	services appservices.Services
}

func NewMarketingHandler(services appservices.Services) *MarketingHandler {
	return &MarketingHandler{services: services}
}

func (r *MarketingHandler) MarketingRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	marketingGroup := rg.Group("/marketing")
	{
		marketingGroup.Use(m.AuthJwtTokenMiddleware())
		bannersGroup := marketingGroup.Group("/banners")
		bannersGroup.POST("", m.RBACPermission("marketing:create"), r.CreateBannerHandler)
		bannersGroup.GET("", m.RBACPermission("marketing:list"), r.GetBannersHandler)
		bannersGroup.GET("/:id", m.RBACPermission("marketing:get"), r.GetBannerByIDHandler)
		bannersGroup.PUT("/:id", m.RBACPermission("marketing:update"), r.UpdateBannerHandler)
		bannersGroup.DELETE("/:id", m.RBACPermission("marketing:delete"), r.DeleteBannerHandler)
	}
}
