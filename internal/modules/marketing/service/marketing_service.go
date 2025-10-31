package marketingservice

import "github.com/vitalfit/api/internal/store"

type MarketingService struct {
	store store.Storage
}

func NewMarketingService(store store.Storage) *MarketingService {
	return &MarketingService{
		store: store,
	}
}
