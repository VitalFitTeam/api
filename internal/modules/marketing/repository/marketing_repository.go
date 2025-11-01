package marketingrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"gorm.io/gorm"
)

type MarketingStore struct {
	db *gorm.DB
}

func NewMarketingStore(db *gorm.DB) *MarketingStore {
	return &MarketingStore{
		db: db,
	}
}

func (s *MarketingStore) CreateBanner(ctx context.Context, banner *marketingdomain.Banner) error {
	err := s.db.WithContext(ctx).Create(banner).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return shared_errors.ErrConflict
			}
		}
		return err
	}
	return nil
}

func (s *MarketingStore) UpdateBanner(ctx context.Context, banner *marketingdomain.Banner) error {
	err := s.db.Save(banner).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return shared_errors.ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *MarketingStore) DeleteBanner(ctx context.Context, bannerID uuid.UUID) error {
	err := s.db.Delete(&marketingdomain.Banner{}, bannerID).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return shared_errors.ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *MarketingStore) GetBannerByID(ctx context.Context, bannerID uuid.UUID) (*marketingdomain.Banner, error) {
	banner := &marketingdomain.Banner{}
	err := s.db.First(banner, bannerID).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return nil, shared_errors.ErrNotFound
		default:
			return nil, err
		}
	}
	return banner, nil
}

func (s *MarketingStore) GetBanners(ctx context.Context) ([]*marketingdomain.Banner, error) {
	banners := []*marketingdomain.Banner{}
	err := s.db.Find(&banners).Error
	if err != nil {
		return nil, err
	}
	return banners, nil
}
