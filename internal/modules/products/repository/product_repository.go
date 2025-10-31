package productsrepository

import "gorm.io/gorm"

type ProductsStore struct {
	db *gorm.DB
}

func NewProductsStore(db *gorm.DB) *ProductsStore {
	return &ProductsStore{
		db: db,
	}
}
