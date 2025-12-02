package reporthandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type ReportHandlersInterface interface {
	ReportRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
}

type ReportHanlders struct {
	services appservices.Services
}

func NewReportHandlers(services appservices.Services) *ReportHanlders {
	return &ReportHanlders{
		services: services,
	}
}

func (r *ReportHanlders) ReportRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	billingGroup := rg.Group("/billing")
	billingGroup.Use(m.AuthJwtTokenMiddleware())
	billingGroup.GET("/stats/global", m.RBACPermission("billing:list"), r.GetGlobalSalesStatsHandler)

}
