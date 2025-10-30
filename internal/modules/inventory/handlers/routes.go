package inventoryhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

func RegisterInventoryRoutes(router *gin.Engine, services appservices.Services, m *auth.AuthMiddleware) {
	v1 := router.Group("/v1")
	h := NewInventoryHandlers(services)
	h.InventoryRoutes(v1, m)
}

func (r *InventoryHandlers) InventoryRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	// ==============================
	//  EQUIPMENT CATALOG
	// ==============================
	equipmentGroup := rg.Group("/equipment-types").
		Use(m.AuthJwtTokenMiddleware())

	{
		equipmentGroup.POST("", r.CreateEquipmentHandler)
		equipmentGroup.GET("", r.GetEquipmentsHandler)
		equipmentGroup.PUT("/:id", r.UpdateEquipmentHandler)
		equipmentGroup.DELETE("/:id", r.DeleteEquipmentHandler)
	}

	// ==============================
	//  BRANCH INVENTORY
	// ==============================
	branchInventoryGroup := rg.Group("/inventory/branch/:branchId").
		Use(m.AuthJwtTokenMiddleware())

	{
		branchInventoryGroup.POST("", r.AddInventoryItemHandler)
		branchInventoryGroup.GET("", r.ListBranchInventoryHandler)
		branchInventoryGroup.PATCH("/:inventoryId", r.UpdateInventoryItemHandler)
		branchInventoryGroup.DELETE("/:inventoryId", r.DeleteInventoryItemHandler)
	}
}
