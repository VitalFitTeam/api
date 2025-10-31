package apphandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
	authhandlers "github.com/vitalfit/api/internal/modules/auth/handlers"
	branchhandlers "github.com/vitalfit/api/internal/modules/branches/handlers"
	instructorhandler "github.com/vitalfit/api/internal/modules/instructor/handler"
	inventoryhandlers "github.com/vitalfit/api/internal/modules/inventory/handlers"
)

type Handlers struct {
	AuthHandlers       authhandlers.AuthHandlersInterface
	BranchHandlers     branchhandlers.BranchHandlersInterface
	InventoryHandlers  inventoryhandlers.InventoryHandlersInterface
	InstructorHandlers instructorhandler.InstructorHandlersInterface
}

func NewAppHandlers(services appservices.Services) Handlers {
	return Handlers{
		AuthHandlers:       authhandlers.NewAuthHandlers(services),
		BranchHandlers:     branchhandlers.NewBranchHandlers(services),
		InventoryHandlers:  inventoryhandlers.NewInventoryHandlers(services),
		InstructorHandlers: instructorhandler.NewInstructorHandlers(services),
	}
}
