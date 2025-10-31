package marketingrepository

import "gorm.io/gorm"

type MarketingStore struct {
	db *gorm.DB
}

func NewMarketingStore(db *gorm.DB) *MarketingStore {
	return &MarketingStore{
		db: db,
	}
}
