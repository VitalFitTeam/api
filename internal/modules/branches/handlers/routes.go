package branchhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

func (r *BranchHandlers) BranchRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {

	branchGroup := rg.Group("/branches").
		Use(m.AuthJwtTokenMiddleware())

	{

		branchGroup.GET("", m.RBACPermission("branches:list"), r.GetBranchesHandler)
		branchGroup.GET("/:id", m.RBACPermission("branches:get"), r.GetBranchByIDHandler)
		branchGroup.POST("", m.RBACPermission("branches:create"), r.CreateBranchHandler)
		branchGroup.PUT("/:id", m.RBACPermission("branches:update"), r.UpdateBranchHandler)
		branchGroup.DELETE("/:id", m.RBACPermission("branches:delete"), r.DeleteBranchHandler)

		branchGroup.GET("/payment-methods", m.RBACPermission("branches:list"), r.GetPaymentMethodsHandler)
		branchGroup.GET("/status", m.RBACPermission("branches:list"), r.GetBranchStatusCount)
	}
}
