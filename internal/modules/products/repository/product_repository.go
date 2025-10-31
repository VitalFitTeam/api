package productsrepository

import (
	"context"

	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	"github.com/vitalfit/api/pkg/db"
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

func (s *ProductsStore) GetServicesCategories(ctx context.Context) ([]productsdomain.ServiceCategory, error) {
	var serviceCategories []productsdomain.ServiceCategory

	err := db.WithTX(s.db, func(tx *gorm.DB) error {
		var listErr error

		serviceCategories, listErr = s.ListServiceCategories(ctx, tx)
		if listErr != nil {
			return listErr
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return serviceCategories, nil
}
