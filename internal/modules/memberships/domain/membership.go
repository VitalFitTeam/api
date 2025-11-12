package membershipsdomain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipType struct {
	MembershipTypeID uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"membership_type_id"`
	Name             string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	DurationDays     int            `gorm:"not null" json:"duration_days"`
	Price            float64        `gorm:"type:numeric(10,2);not null" json:"price"`
	IsActive         bool           `gorm:"not null;default:true" json:"is_active"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (MembershipType) TableName() string {
	return "membership_types"
}

type MembershipSummary struct {
	Total     int64 `json:"total"`
	Actives   int64 `json:"actives"`
	Inactives int64 `json:"inactives"`
}
