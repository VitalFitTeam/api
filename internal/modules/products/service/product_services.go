package productsservice

import "github.com/vitalfit/api/internal/store"

type ProductsService struct {
	store store.Storage
}

func NewProductsService(store store.Storage) *ProductsService {
	return &ProductsService{
		store: store,
	}
}
