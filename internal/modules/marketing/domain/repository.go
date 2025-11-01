package marketingdomain

import (
	"context"

	"github.com/google/uuid"
)

type MarketingRepository interface {
	CreateBanner(context.Context, *Banner) error
	UpdateBanner(context.Context, *Banner) error
	DeleteBanner(context.Context, uuid.UUID) error
	GetBannerByID(context.Context, uuid.UUID) (*Banner, error)
	GetBanners(context.Context) ([]*Banner, error)
}
