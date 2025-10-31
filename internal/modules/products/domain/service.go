package productsdomain

import (
	"context"
)

type ProductsServiceInterface interface {
	ListServiceCategories(ctx context.Context) ([]ServiceCategory, error)
}
