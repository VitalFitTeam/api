package staffhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type StaffHandlersInterface interface {
	StaffRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)
	AssignStaffToBranchHandler(c *gin.Context)
	ListBranchStaffByRoleHandler(c *gin.Context)
	RemoveStaffFromBranchHandler(c *gin.Context)
}

type StaffHandlers struct {
	services appservices.Services
}

func NewStaffHandlers(services appservices.Services) *StaffHandlers {
	return &StaffHandlers{services: services}
}

func (h *StaffHandlers) StaffRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	staffRoutes := rg.Group("/branches/:id/staff")
	{
		staffRoutes.Use(m.AuthJwtTokenMiddleware())
		staffRoutes.POST("", m.RBACPermission("users:update"), h.AssignStaffToBranchHandler)
		staffRoutes.GET("", m.RBACPermission("users:list"), h.ListBranchStaffByRoleHandler)
		staffRoutes.DELETE("/:staffId", m.RBACPermission("users:update"), h.RemoveStaffFromBranchHandler)
	}
}
