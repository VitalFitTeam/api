package marketingdomain

import (
	"github.com/google/uuid"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
)

type Banner struct {
	BannerID uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"banner_id"`
	Name     string    `gorm:"type:varchar(255);not null" json:"name"`
	ImageURL string    `gorm:"type:varchar(255);not null" json:"image_url"`
	LinkURL  string    `gorm:"type:varchar(255)" json:"link_url"`
	IsActive bool      `gorm:"not null;default:true" json:"is_active"`

	Services []productsdomain.Service `gorm:"many2many:banner_services;foreignKey:BannerID;joinForeignKey:BannerID;References:ServiceID;joinReferences:ServiceID" json:"services,omitempty"`
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
