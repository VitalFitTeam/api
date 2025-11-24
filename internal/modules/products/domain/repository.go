package productsdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
)

type ProductsRepository interface {
	CreateServiceCategory(ctx context.Context, tx *gorm.DB, serviceCategory *ServiceCategory) error
	ListServiceCategories(ctx context.Context, tx *gorm.DB) ([]ServiceCategory, error)
	GetServicesCategories(ctx context.Context) ([]ServiceCategory, error)

	CreateService(ctx context.Context, service *Service, bannerID uuid.UUID) error
	CreateServiceTX(ctx context.Context, tx *gorm.DB, service *Service, bannerID uuid.UUID) error
	GetServices(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]Service, error)
	GetTotalCount(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error)
	GetServiceSummary(ctx context.Context) (*ServicesSummary, error)
	DeleteService(ctx context.Context, serviceID uuid.UUID) error
	GetServiceByID(ctx context.Context, serviceID uuid.UUID) (*Service, error)
	UpdateService(ctx context.Context, service *Service, bannerID uuid.UUID) error

	AssignBranchService(ctx context.Context, branchServices []*ServiceBranchDetail) error
	GetBranchService(ctx context.Context, branchID uuid.UUID) ([]*ServiceBranchDetail, error)
	UpdateBranchService(ctx context.Context, branchService *ServiceBranchDetail) error
	DeleteBranchService(ctx context.Context, branchID uuid.UUID, serviceID uuid.UUID) error
	GetBranchServiceByID(ctx context.Context, branchID uuid.UUID, serviceID uuid.UUID) (*ServiceBranchDetail, error)

	GetPublicServices(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]ServiceWithPrice, int64, error)
}
