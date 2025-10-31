package productsservice

import (
	"context"

	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"github.com/vitalfit/api/internal/store"
)

type ProductsService struct {
	store store.Storage
}

func NewProductsService(store store.Storage) *ProductsService {
	return &ProductsService{
		store: store,
	}
}

func (s *ProductsService) ListServiceCategories(ctx context.Context) ([]productsdomain.ServiceCategory, error) {
	return s.store.Products.GetServicesCategories(ctx)
}
