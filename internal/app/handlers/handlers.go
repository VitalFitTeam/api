package apphandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
	authhandlers "github.com/vitalfit/api/internal/modules/auth/handlers"
	branchhandlers "github.com/vitalfit/api/internal/modules/branches/handlers"
	instructorhandler "github.com/vitalfit/api/internal/modules/instructor/handler"
	inventoryhandlers "github.com/vitalfit/api/internal/modules/inventory/handlers"
	marketinghandlers "github.com/vitalfit/api/internal/modules/marketing/handlers"
	productshandler "github.com/vitalfit/api/internal/modules/products/handler"
)

type Handlers struct {
	AuthHandlers       authhandlers.AuthHandlersInterface
	BranchHandlers     branchhandlers.BranchHandlersInterface
	InventoryHandlers  inventoryhandlers.InventoryHandlersInterface
	InstructorHandlers instructorhandler.InstructorHandlersInterface
	ProductsHandlers   productshandler.ProductsHandlerInterface
	MarketingHandlers  marketinghandlers.MarketingHandlerInterface
}

func NewAppHandlers(services appservices.Services) Handlers {
	return Handlers{
		AuthHandlers:       authhandlers.NewAuthHandlers(services),
		BranchHandlers:     branchhandlers.NewBranchHandlers(services),
		InventoryHandlers:  inventoryhandlers.NewInventoryHandlers(services),
		InstructorHandlers: instructorhandler.NewInstructorHandlers(services),
		ProductsHandlers:   productshandler.NewProductsHandler(services),
		MarketingHandlers:  marketinghandlers.NewMarketingHandler(services),
	}
}
