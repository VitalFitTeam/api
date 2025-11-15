package store

import (
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	authrepository "github.com/vitalfit/api/internal/modules/auth/repository"
	billingdomain "github.com/vitalfit/api/internal/modules/billing/domain"
	billingrepository "github.com/vitalfit/api/internal/modules/billing/repository"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	branchrepository "github.com/vitalfit/api/internal/modules/branches/repository"
	combosdomain "github.com/vitalfit/api/internal/modules/combos/domain"
	combosrepository "github.com/vitalfit/api/internal/modules/combos/repository"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	instructorrepository "github.com/vitalfit/api/internal/modules/instructor/repository"
	inventorydomain "github.com/vitalfit/api/internal/modules/inventory/domain"
	inventoryrepository "github.com/vitalfit/api/internal/modules/inventory/repository"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
	marketingrepository "github.com/vitalfit/api/internal/modules/marketing/repository"
	membershipsdomain "github.com/vitalfit/api/internal/modules/memberships/domain"
	membershipsrepository "github.com/vitalfit/api/internal/modules/memberships/repository"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	productsrepository "github.com/vitalfit/api/internal/modules/products/repository"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	schedulerepository "github.com/vitalfit/api/internal/modules/schedule/repository"
	"gorm.io/gorm"
)

type Storage struct {
	Users           authdomain.UserRepository
	Roles           authdomain.RolesRepository
	Branches        branchdomain.BranchesRepository
	Locations       branchdomain.LocationRepository
	PaymentMethods  billingdomain.PaymentMethodsRepository
	Equipment       inventorydomain.EquipmentRepository
	BranchInventory inventorydomain.BranchInventoryRepository
	Instructor      instructordomain.InstructorRepository
	Products        productsdomain.ProductsRepository
	Marketing       marketingdomain.MarketingRepository
	Membership      membershipsdomain.MembershipsRepository
	Billing         billingdomain.BillingRepository
	Schedule        scheduledomain.ScheduleRepository
	Combos          combosdomain.CombosRepository
	FiscalDocuments billingdomain.FiscalDocumentRepository
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Users:           authrepository.NewUserStore(db),
		Roles:           authrepository.NewRoleStore(db),
		Branches:        branchrepository.NewBranchesStore(db),
		Locations:       branchrepository.NewLocationsStore(db),
		PaymentMethods:  billingrepository.NewPaymentMethodsStore(db),
		Equipment:       inventoryrepository.NewEquipmentStore(db),
		BranchInventory: inventoryrepository.NewBranchInventoryStore(db),
		Instructor:      instructorrepository.NewInstructorStore(db),
		Products:        productsrepository.NewProductsStore(db),
		Marketing:       marketingrepository.NewMarketingStore(db),
		Membership:      membershipsrepository.NewMembershipStore(db),
		Billing:         billingrepository.NewBillingStore(db),
		Schedule:        schedulerepository.NewScheduleStore(db),
		Combos:          combosrepository.NewCombosStore(db),
		FiscalDocuments: billingrepository.NewFiscalDocumentStore(db),
	}
}
