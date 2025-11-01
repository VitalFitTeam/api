package marketingdomain

import (
	"github.com/google/uuid"
)

type Banner struct {
	BannerID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"banner_id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	ImageURL string    `gorm:"type:varchar(255);not null" json:"image_url"`
	LinkURL  string    `gorm:"type:varchar(255)" json:"link_url"`
	IsActive bool      `gorm:"not null;default:true" json:"is_active"`
}

func (Banner) TableName() string {
	return "banners"
}

type BannerService struct {
	BannerID  uuid.UUID `gorm:"primaryKey"`
	ServiceID uuid.UUID `gorm:"primaryKey"`
}

func (BannerService) TableName() string {
	return "banner_services"
}
