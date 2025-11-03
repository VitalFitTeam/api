package membershipsdomain

import (
	"time"

	"github.com/google/uuid"
)

type MembershipType struct {
	MembershipTypeID uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"membership_type_id"`
	Name             string     `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description      string     `gorm:"type:text" json:"description"`
	DurationDays     int        `gorm:"not null" json:"duration_days"`
	Price            float64    `gorm:"type:numeric(10,2);not null" json:"price"`
	IsActive         bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (MembershipType) TableName() string {
	return "membership_types"
}
