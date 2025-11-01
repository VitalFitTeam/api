package marketinghandlers

import (
	"github.com/google/uuid"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
)

type CreateBannerPayload struct {
	Name     string `json:"name" binding:"required"`
	ImageURL string `json:"image_url" binding:"required"`
	LinkURL  string `json:"link_url" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type UpdateBannerPayload struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url" `
	IsActive bool   `json:"is_active"`
}

func (s *CreateBannerPayload) toBanner() (*marketingdomain.Banner, error) {
	banner := &marketingdomain.Banner{
		Name:     s.Name,
		ImageURL: s.ImageURL,
		LinkURL:  s.LinkURL,
		IsActive: s.IsActive,
	}
	return banner, nil
}

func (s *UpdateBannerPayload) toBanner() (*marketingdomain.Banner, error) {
	banner := &marketingdomain.Banner{
		Name:     s.Name,
		ImageURL: s.ImageURL,
		LinkURL:  s.LinkURL,
		IsActive: s.IsActive,
	}
	return banner, nil
}

type BannerResponse struct {
	BannerID uuid.UUID `json:"banner_id"`
	Name     string    `json:"name"`
	ImageURL string    `json:"image_url"`
	LinkURL  string    `json:"link_url"`
	IsActive bool      `json:"is_active"`
}
