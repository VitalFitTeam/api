package apphandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
	authhandlers "github.com/vitalfit/api/internal/modules/auth/handlers"
	branchhandlers "github.com/vitalfit/api/internal/modules/branches/handlers"
	inventoryhandlers "github.com/vitalfit/api/internal/modules/inventory/handlers"
)

type Handlers struct {
	AuthHandlers      authhandlers.AuthHandlersInterface
	BranchHandlers    branchhandlers.BranchHandlersInterface
	InventoryHandlers inventoryhandlers.InventoryHandlersInterface
}

func NewAppHandlers(services appservices.Services) Handlers {
	inventoryHandler := inventoryhandlers.NewInventoryHandlers(services)
	branchHandler := branchhandlers.NewBranchHandlers(services)
	branchHandler.SetInventoryHandlers(inventoryHandler)
	return Handlers{
		AuthHandlers:      authhandlers.NewAuthHandlers(services),
		BranchHandlers:    branchHandler,
		InventoryHandlers: inventoryHandler,
	}

}
