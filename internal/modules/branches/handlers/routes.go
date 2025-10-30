package branchhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

func (r *BranchHandlers) BranchRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	branchGroup := rg.Group("/branches").
		Use(m.AuthJwtTokenMiddleware(), m.CheckRoleAccess("super_admin"))

	{
		branchGroup.GET("", r.GetBranchesHandler)
		branchGroup.GET("/:id", r.GetBranchByIDHandler)
		branchGroup.POST("", r.CreateBranchHandler)
		branchGroup.GET("/payment-methods", r.GetPaymentMethodsHandler)
		branchGroup.PUT("/:id", r.UpdateBranchHandler)
		branchGroup.DELETE("/:id", r.DeleteBranchHandler)
		branchGroup.GET("/status", r.GetBranchStatusCount)
	}

	if r.inventoryHandlers != nil {
		r.inventoryHandlers.InventoryRoutes(rg, m)
	}
}
