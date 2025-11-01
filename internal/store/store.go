package store

import (
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	authrepository "github.com/vitalfit/api/internal/modules/auth/repository"
	branchdomain "github.com/vitalfit/api/internal/modules/branches/domain"
	branchrepository "github.com/vitalfit/api/internal/modules/branches/repository"
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
	"gorm.io/gorm"
)

type Storage struct {
	Users           authdomain.UserRepository
	Roles           authdomain.RolesRepository
	Branches        branchdomain.BranchesRepository
	Locations       branchdomain.LocationRepository
	PaymentMethods  branchdomain.PaymentMethodsRepository
	Equipment       inventorydomain.EquipmentRepository
	BranchInventory inventorydomain.BranchInventoryRepository
	Instructor      instructordomain.InstructorRepository
	Products        productsdomain.ProductsRepository
	Marketing       marketingdomain.MarketingRepository
	Membership      membershipsdomain.MembershipsRepository
}

func NewStorage(db *gorm.DB) Storage {
	return Storage{
		Users:           authrepository.NewUserStore(db),
		Roles:           authrepository.NewRoleStore(db),
		Branches:        branchrepository.NewBranchesStore(db),
		Locations:       branchrepository.NewLocationsStore(db),
		PaymentMethods:  branchrepository.NewPaymentMethodsStore(db),
		Equipment:       inventoryrepository.NewEquipmentStore(db),
		BranchInventory: inventoryrepository.NewBranchInventoryStore(db),
		Instructor:      instructorrepository.NewInstructorStore(db),
		Products:        productsrepository.NewProductsStore(db),
		Marketing:       marketingrepository.NewMarketingStore(db),
		Membership:      membershipsrepository.NewMembershipStore(db),
	}
}
