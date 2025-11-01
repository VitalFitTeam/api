package productsdomain

import (
	"context"

	"github.com/google/uuid"
)

type ProductsServiceInterface interface {
	ListServiceCategories(ctx context.Context) ([]ServiceCategory, error)
	CreateService(ctx context.Context, service *Service, bannerID uuid.UUID) error
	GetServices(ctx context.Context) ([]Service, error)
	DeleteService(ctx context.Context, serviceID uuid.UUID) error
	GetServiceByID(ctx context.Context, serviceID uuid.UUID) (*Service, error)
	UpdateService(ctx context.Context, service *Service, bannerID uuid.UUID) error
}
