package branchhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

func (r *BranchHandlers) BranchRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	branchGroup := rg.Group("/branches").Use(m.AuthJwtTokenMiddleware())
	{
		branchGroup.GET("/", r.GetBranchesHandler).Use(m.CheckRoleAccess("super_admin"))
		branchGroup.POST("/", r.CreateBranchHandler).Use(m.CheckRoleAccess("super_admin"))
		branchGroup.GET("/payment-methods", r.GetPaymentMethodsHandler).Use(m.CheckRoleAccess("super_admin"))
	}
}
