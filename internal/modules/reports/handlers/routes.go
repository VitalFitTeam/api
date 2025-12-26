package reporthandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type ReportHandlersInterface interface {
	ReportRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	GetGlobalSalesStatsHandler(c *gin.Context)
	GetMonthlySalesKPIHandler(c *gin.Context)
	GetActiveMembersKPIHandler(c *gin.Context)
	GetOccupancyKPIHandler(c *gin.Context)
	GetTotalSalesHandler(c *gin.Context)
	GetTopBranchesPerformanceHandler(c *gin.Context)
	GetTotalClientsStatHandler(c *gin.Context)
	GetActiveBranchesCountHandler(c *gin.Context)
	GetSalesByCategoryHandler(c *gin.Context)
	GetTopInstructorsByAttendanceHandler(c *gin.Context)
	GetSalesByPaymentMethodHandler(c *gin.Context)
	GetSalesByHourHandler(c *gin.Context)
	GetMostUsedServicesHandler(c *gin.Context)
	GetWeeklySalesChartHandler(c *gin.Context)
	GetActivityHeatmapHandler(c *gin.Context)
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
	billingGroup := rg.Group("/reports")
	billingGroup.Use(m.AuthJwtTokenMiddleware(), m.RBACPermission("reports:view_sales"))

	statsGroup := billingGroup.Group("/stats")
	statsGroup.GET("/global", r.GetGlobalSalesStatsHandler)
	statsGroup.GET("/total-sales", r.GetTotalSalesHandler)
	statsGroup.GET("/top-branches", r.GetTopBranchesPerformanceHandler)
	statsGroup.GET("/total-clients", r.GetTotalClientsStatHandler)
	statsGroup.GET("/total-active-branches", r.GetActiveBranchesCountHandler)

	kpiGroup := billingGroup.Group("/kpi")
	kpiGroup.GET("/monthly-sales", r.GetMonthlySalesKPIHandler)
	kpiGroup.GET("/active-members", r.GetActiveMembersKPIHandler)
	kpiGroup.GET("/occupancy", r.GetOccupancyKPIHandler)

	chartsGroup := billingGroup.Group("/charts")
	chartsGroup.GET("/sales-by-category", r.GetSalesByCategoryHandler)
	chartsGroup.GET("/top-instructors", r.GetTopInstructorsByAttendanceHandler)
	chartsGroup.GET("/sales-by-payment-method", r.GetSalesByPaymentMethodHandler)
	chartsGroup.GET("/sales-by-hour", r.GetSalesByHourHandler)
	chartsGroup.GET("/most-used-services", r.GetMostUsedServicesHandler)
	chartsGroup.GET("/weekly-sales", r.GetWeeklySalesChartHandler)
	chartsGroup.GET("/activity-heatmap", r.GetActivityHeatmapHandler)

}
