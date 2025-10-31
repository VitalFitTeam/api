package inventoryhandlers

import (
	"github.com/gin-gonic/gin"
	appservices "github.com/vitalfit/api/internal/app/services"
	"github.com/vitalfit/api/internal/shared/middleware/auth"
)

type InventoryHandlersInterface interface {
	InventoryRoutes(rg *gin.RouterGroup, m *auth.AuthMiddleware)

	// Equipment
	CreateEquipmentHandler(c *gin.Context)
	GetEquipmentsHandler(c *gin.Context)
	UpdateEquipmentHandler(c *gin.Context)
	DeleteEquipmentHandler(c *gin.Context)

	// Branch inventory
	AddInventoryItemHandler(c *gin.Context)
	ListBranchInventoryHandler(c *gin.Context)
	UpdateInventoryItemHandler(c *gin.Context)
	DeleteInventoryItemHandler(c *gin.Context)
}

type InventoryHandlers struct {
	services appservices.Services
}

// Constructor
func NewInventoryHandlers(services appservices.Services) *InventoryHandlers {
	return &InventoryHandlers{services: services}
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
	branchInventoryGroup := rg.Group("/branches/:id/equipment").
		Use(m.AuthJwtTokenMiddleware())

	{
		branchInventoryGroup.POST("", r.AddInventoryItemHandler)
		branchInventoryGroup.GET("", r.ListBranchInventoryHandler)
		branchInventoryGroup.PATCH("/:inventoryId", r.UpdateInventoryItemHandler)
		branchInventoryGroup.DELETE("/:inventoryId", r.DeleteInventoryItemHandler)
	}
}
