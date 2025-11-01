package productsdomain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductsRepository interface {
	CreateServiceCategory(ctx context.Context, tx *gorm.DB, serviceCategory *ServiceCategory) error
	ListServiceCategories(ctx context.Context, tx *gorm.DB) ([]ServiceCategory, error)
	GetServicesCategories(ctx context.Context) ([]ServiceCategory, error)

	CreateService(ctx context.Context, service *Service, bannerID uuid.UUID) error
	CreateServiceTX(ctx context.Context, tx *gorm.DB, service *Service, bannerID uuid.UUID) error
	GetServices(ctx context.Context) ([]Service, error)
	DeleteService(ctx context.Context, serviceID uuid.UUID) error
	GetServiceByID(ctx context.Context, serviceID uuid.UUID) (*Service, error)
	UpdateService(ctx context.Context, service *Service, bannerID uuid.UUID) error

	AssignBranchService(ctx context.Context, branchServices []*ServiceBranchDetail) error
	GetBranchService(ctx context.Context, branchID uuid.UUID) ([]*ServiceBranchDetail, error)
	UpdateBranchService(ctx context.Context, branchService *ServiceBranchDetail) error
	DeleteBranchService(ctx context.Context, branchID uuid.UUID, serviceID uuid.UUID) error
	GetBranchServiceByID(ctx context.Context, branchID uuid.UUID, serviceID uuid.UUID) (*ServiceBranchDetail, error)
}
