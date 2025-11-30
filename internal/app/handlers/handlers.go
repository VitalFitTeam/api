package apphandlers

import (
	appservices "github.com/vitalfit/api/internal/app/services"
	accesshandler "github.com/vitalfit/api/internal/modules/access/handler"
	authhandlers "github.com/vitalfit/api/internal/modules/auth/handlers"
	billinghandlers "github.com/vitalfit/api/internal/modules/billing/handlers"
	bookinghandlers "github.com/vitalfit/api/internal/modules/booking/handlers"
	branchhandlers "github.com/vitalfit/api/internal/modules/branches/handlers"
	comboshandler "github.com/vitalfit/api/internal/modules/combos/handler"
	instructorhandler "github.com/vitalfit/api/internal/modules/instructor/handler"
	inventoryhandlers "github.com/vitalfit/api/internal/modules/inventory/handlers"
	marketinghandlers "github.com/vitalfit/api/internal/modules/marketing/handlers"
	membershipshandlers "github.com/vitalfit/api/internal/modules/memberships/handlers"
	productshandler "github.com/vitalfit/api/internal/modules/products/handler"
	schedulehandlers "github.com/vitalfit/api/internal/modules/schedule/handlers"
)

type Handlers struct {
	AuthHandlers       authhandlers.AuthHandlersInterface
	BranchHandlers     branchhandlers.BranchHandlersInterface
	InventoryHandlers  inventoryhandlers.InventoryHandlersInterface
	InstructorHandlers instructorhandler.InstructorHandlersInterface
	ProductsHandlers   productshandler.ProductsHandlerInterface
	MarketingHandlers  marketinghandlers.MarketingHandlerInterface
	MembershipHandlers membershipshandlers.MembershipsHandlerInterface
	BillingHandlers    billinghandlers.BillingHandlersInterface
	ScheduleHandlers   schedulehandlers.ScheduleHandlersInterface
	CombosHandlers     comboshandler.CombosHandlerInterface
	BookingHandlers    bookinghandlers.BookingHandlersInterface
	AccessHandlers     accesshandler.AccessHandlerInterface
}

func NewAppHandlers(services appservices.Services) Handlers {
	return Handlers{
		AuthHandlers:       authhandlers.NewAuthHandlers(services),
		BranchHandlers:     branchhandlers.NewBranchHandlers(services),
		InventoryHandlers:  inventoryhandlers.NewInventoryHandlers(services),
		InstructorHandlers: instructorhandler.NewInstructorHandlers(services),
		ProductsHandlers:   productshandler.NewProductsHandler(services),
		MarketingHandlers:  marketinghandlers.NewMarketingHandler(services),
		MembershipHandlers: membershipshandlers.NewMembershipHandler(services),
		BillingHandlers:    billinghandlers.NewBillingHandlers(services),
		ScheduleHandlers:   schedulehandlers.NewScheduleHandlers(services),
		CombosHandlers:     comboshandler.NewCombosHandler(services),
		BookingHandlers:    bookinghandlers.NewBookingHandlers(services),
		AccessHandlers:     accesshandler.NewAccessHandler(services),
	}
}
