package productsservice

<<<<<<< HEAD
import (
	"context"

	"github.com/google/uuid"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"github.com/vitalfit/api/internal/store"
)
=======
import "github.com/vitalfit/api/internal/store"
>>>>>>> dev

type ProductsService struct {
	store store.Storage
}

func NewProductsService(store store.Storage) *ProductsService {
	return &ProductsService{
		store: store,
	}
}
<<<<<<< HEAD

func (s *ProductsService) ListServiceCategories(ctx context.Context) ([]productsdomain.ServiceCategory, error) {
	return s.store.Products.GetServicesCategories(ctx)
}

func (s *ProductsService) CreateService(ctx context.Context, service *productsdomain.Service, bannerID uuid.UUID) error {
	return s.store.Products.CreateService(ctx, service, bannerID)
}

func (s *ProductsService) GetServices(ctx context.Context) ([]productsdomain.Service, error) {
	return s.store.Products.GetServices(ctx)
}

func (s *ProductsService) DeleteService(ctx context.Context, serviceID uuid.UUID) error {
	return s.store.Products.DeleteService(ctx, serviceID)
}

func (s *ProductsService) GetServiceByID(ctx context.Context, serviceID uuid.UUID) (*productsdomain.Service, error) {
	return s.store.Products.GetServiceByID(ctx, serviceID)
}

func (s *ProductsService) UpdateService(ctx context.Context, service *productsdomain.Service, bannerID uuid.UUID) error {
	return s.store.Products.UpdateService(ctx, service, bannerID)
}
=======
>>>>>>> dev
