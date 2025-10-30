package inventoryhandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

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
	branchInventoryGroup := rg.Group("/branches/:id/equipment").
		Use(m.AuthJwtTokenMiddleware())

	{
		branchInventoryGroup.POST("", r.AddInventoryItemHandler)
		branchInventoryGroup.GET("", r.ListBranchInventoryHandler)
		branchInventoryGroup.PATCH("/:inventoryId", r.UpdateInventoryItemHandler)
		branchInventoryGroup.DELETE("/:inventoryId", r.DeleteInventoryItemHandler)
	}
}
