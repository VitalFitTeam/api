package productsdomain

import (
	"time"

	"github.com/google/uuid"
	marketingdomain "github.com/vitalfit/api/internal/modules/marketing/domain"
	"gorm.io/gorm"
)

type ServiceCategory struct {
	CategoryID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name       string    `gorm:"type:varchar(100);unique;not null"`

	Services []Service `gorm:"foreignKey:CategoryID"`
}

func (ServiceCategory) TableName() string {
	return "service_categories"
}

type Service struct {
	ServiceID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null"`
	Name            string    `gorm:"type:varchar(255);unique;not null"`
	Description     string    `gorm:"type:text"`
	DurationMinutes int64     `gorm:"type:int"`
	PriorityScore   int64     `gorm:"type:int;default:50"`
	IsFeatured      bool      `gorm:"type:boolean;default:false;not null"`
	CreatedAt       time.Time `gorm:"default:now()"`
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Category ServiceCategory `gorm:"foreignKey:CategoryID;references:CategoryID"`

	Images []ServiceImage `gorm:"foreignKey:ServiceID"`

	BranchDetails []ServiceBranchDetail    `gorm:"foreignKey:ServiceID"`
	Banners       []marketingdomain.Banner `gorm:"many2many:banner_services;foreignKey:ServiceID;joinForeignKey:ServiceID;References:BannerID;joinReferences:BannerID" json:"banners,omitempty"`
}

func (Service) TableName() string {
	return "services"
}

type ServiceImage struct {
	ImageID      uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ServiceID    uuid.UUID `gorm:"type:uuid;not null"`
	ImageURL     string    `gorm:"type:varchar(255);not null"`
	AltText      string    `gorm:"type:varchar(255)"`
	DisplayOrder int       `gorm:"type:int;not null;default:0"`
	IsPrimary    bool      `gorm:"type:boolean;not null;default:false"`

	Service Service `gorm:"foreignKey:ServiceID"`
}

func (ServiceImage) TableName() string {
	return "service_images"
}

type ServiceBranchDetail struct {
	ServiceID         uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	BranchID          uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	IsVisible         bool      `gorm:"type:boolean;not null;default:true"`
	MaxCapacity       int       `gorm:"type:int;not null"`
	PriceForMember    float64   `gorm:"type:decimal(10,2);not null"`
	PriceForNonMember float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt         time.Time `gorm:"default:now()"`
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	Service Service `gorm:"foreignKey:ServiceID"`
}

func (ServiceBranchDetail) TableName() string {
	return "service_branch_details"
}
