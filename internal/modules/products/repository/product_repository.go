package productsrepository

import (
	"context"

	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"gorm.io/gorm"
)

type ProductsStore struct {
	db *gorm.DB
}

func NewProductsStore(db *gorm.DB) *ProductsStore {
	return &ProductsStore{
		db: db,
	}
}

func (s *ProductsStore) CreateServiceCategory(ctx context.Context, tx *gorm.DB, serviceCategory *productsdomain.ServiceCategory) error {
	if err := tx.Create(serviceCategory).Error; err != nil {
		return err
	}
	return nil
}

func (s *ProductsStore) ListServiceCategories(ctx context.Context, tx *gorm.DB) ([]productsdomain.ServiceCategory, error) {
	var serviceCategories []productsdomain.ServiceCategory
	if err := tx.Find(&serviceCategories).Error; err != nil {
		return nil, err
	}
	return serviceCategories, nil

}
