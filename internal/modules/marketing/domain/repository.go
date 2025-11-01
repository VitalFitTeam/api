package marketingdomain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketingRepository interface {
	CreateBanner(context.Context, *Banner) error
	CreateBannerTX(context.Context, *gorm.DB, *Banner) error
	UpdateBanner(context.Context, *Banner) error
	DeleteBanner(context.Context, uuid.UUID) error
	GetBannerByID(context.Context, uuid.UUID) (*Banner, error)
	GetBanners(context.Context) ([]*Banner, error)
}
