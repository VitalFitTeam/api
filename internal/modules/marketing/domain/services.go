package marketingdomain

import (
	"context"

	"github.com/google/uuid"
)

type MarketingServiceInterface interface {
	// Banner operations
	CreateBanner(context.Context, *Banner) error
	UpdateBanner(context.Context, *Banner) error
	DeleteBanner(context.Context, uuid.UUID) error
	GetBannerByID(context.Context, uuid.UUID) (*Banner, error)
	GetBanners(context.Context) ([]*Banner, error)

	// Promotion operations
	CreatePromotion(context.Context, *Promotion) error
	UpdatePromotion(context.Context, *Promotion) error
	DeletePromotion(context.Context, uuid.UUID) error
	GetPromotionByID(context.Context, uuid.UUID) (*Promotion, error)
	GetPromotions(context.Context) ([]*Promotion, error)
}
