package bookingdomain

import (
	"time"

	"github.com/google/uuid"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	"gorm.io/gorm"
)

type Booking struct {
	BookingID uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"booking_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	ClassID   uuid.UUID      `gorm:"type:uuid;not null" json:"class_id"`
	Status    string         `gorm:"type:varchar(30);not null;default:'Confirmed'" json:"status"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Class scheduledomain.Class `gorm:"foreignKey:ClassID" json:"-"`
}

func (Booking) TableName() string {
	return "bookings"
}
