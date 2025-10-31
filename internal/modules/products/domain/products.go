package productsdomain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceCategory struct {
	CategoryID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name       string    `gorm:"type:varchar(100);unique;not null"`

	Services []Service `gorm:"foreignKey:CategoryID"`
}

type Service struct {
	ServiceID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null"`
	Name            string    `gorm:"type:varchar(255);unique;not null"`
	Description     string    `gorm:"type:text"`
	DurationMinutes int       `gorm:"type:int"`
	PriorityScore   int       `gorm:"type:int;default:50"`
	IsFeatured      bool      `gorm:"type:boolean;default:false;not null"`
	CreatedAt       time.Time `gorm:"default:now()"`
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	ServiceCategory ServiceCategory `gorm:"foreignKey:CategoryID"`

	Images []ServiceImage `gorm:"foreignKey:ServiceID"`

	BranchDetails []ServiceBranchDetail `gorm:"foreignKey:ServiceID"`
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

type ServiceBranchDetail struct {
	ServiceID         uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	BranchID          uuid.UUID `gorm:"type:uuid;not null;primaryKey"` // Asumiendo que tienes un modelo Branch
	IsVisible         bool      `gorm:"type:boolean;not null;default:true"`
	MaxCapacity       int       `gorm:"type:int;not null"`
	PriceForMember    float64   `gorm:"type:decimal(10,2);not null"`
	PriceForNonMember float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt         time.Time `gorm:"default:now()"`
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Relaciones (Belongs To)
	Service Service `gorm:"foreignKey:ServiceID"`
	// Branch  Branch  `gorm:"foreignKey:BranchID"` // Descomenta esto cuando tengas tu struct Branch
}
