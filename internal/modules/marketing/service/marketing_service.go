package marketingservice

import (
	"context"

	"github.com/google/uuid"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
	"github.com/vitalfit/api/internal/store"
)

type MarketingService struct {
	store store.Storage
}

func NewMarketingService(store store.Storage) *MarketingService {
	return &MarketingService{
		store: store,
	}
}

func (s *MarketingService) CreateBanner(ctx context.Context, banner *marketingdomain.Banner) error {
	err := s.store.Marketing.CreateBanner(ctx, banner)
	if err != nil {
		return err
	}
	return nil
}

func (s *MarketingService) UpdateBanner(ctx context.Context, banner *marketingdomain.Banner) error {
	err := s.store.Marketing.UpdateBanner(ctx, banner)
	if err != nil {
		return err
	}
	return nil
}

func (s *MarketingService) DeleteBanner(ctx context.Context, bannerID uuid.UUID) error {
	err := s.store.Marketing.DeleteBanner(ctx, bannerID)
	if err != nil {
		return err
	}
	return nil
}

func (s *MarketingService) GetBannerByID(ctx context.Context, bannerID uuid.UUID) (*marketingdomain.Banner, error) {
	banner, err := s.store.Marketing.GetBannerByID(ctx, bannerID)
	if err != nil {
		return nil, err
	}
	return banner, nil
}

func (s *MarketingService) GetBanners(ctx context.Context) ([]*marketingdomain.Banner, error) {
	banners, err := s.store.Marketing.GetBanners(ctx)
	if err != nil {
		return nil, err
	}
	return banners, nil
}
