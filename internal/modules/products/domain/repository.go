package productsdomain

import (
	"context"

	"gorm.io/gorm"
)

type ProductsRepository interface {
	CreateServiceCategory(ctx context.Context, tx *gorm.DB, serviceCategory *ServiceCategory) error
	ListServiceCategories(ctx context.Context, tx *gorm.DB) ([]ServiceCategory, error)
}
