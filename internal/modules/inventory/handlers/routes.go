package inventoryhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

// ==============================
// INVENTORY ROUTES REGISTRATION
// ==============================

// RegisterInventoryRoutes registra las rutas del módulo de inventario dentro del grupo principal /v1.
// Esto permite que las rutas estén disponibles sin modificar application.go.
func RegisterInventoryRoutes(router *gin.Engine, services appservices.Services, m *auth.AuthMiddleware) {
	v1 := router.Group("/v1")

	h := NewInventoryHandlers(services)
	h.InventoryRoutes(v1, m)
}

func (r *InventoryHandlers) InventoryRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware) {
	// ==============================
	//  EQUIPMENT CATALOG (SUPER ADMIN)
	// ==============================
	equipmentGroup := rg.Group("/equipment-types").
		Use(m.AuthJwtTokenMiddleware(), m.CheckRoleAccess("super_admin"))

	{
		equipmentGroup.POST("", r.CreateEquipmentHandler)
		equipmentGroup.GET("", r.GetEquipmentsHandler)
		equipmentGroup.PUT("/:id", r.UpdateEquipmentHandler)
		equipmentGroup.DELETE("/:id", r.DeleteEquipmentHandler)
	}

	// ==============================
	//  BRANCH INVENTORY (BRANCH ADMIN)
	// ==============================
	branchInventoryGroup := rg.Group("/branches/inventory/:branchId").
		Use(m.AuthJwtTokenMiddleware(), m.CheckRoleAccess("branch_admin"))

	{
		branchInventoryGroup.POST("", r.AddInventoryItemHandler)
		branchInventoryGroup.GET("", r.ListBranchInventoryHandler)
		branchInventoryGroup.PATCH("/:inventoryId", r.UpdateInventoryItemHandler)
		branchInventoryGroup.DELETE("/:inventoryId", r.DeleteInventoryItemHandler)
	}
}
