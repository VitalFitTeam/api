package bookingdomain

import (
	"time"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	productsdomain "github.com/vitalfit/api/internal/modules/products/domain"
	scheduledomain "github.com/vitalfit/api/internal/modules/schedule/domain"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusConfirmed         BookingStatus = "Confirmed"
	BookingStatusCancelledByUser   BookingStatus = "CancelledByUser"
	BookingStatusCancelledBySystem BookingStatus = "CancelledBySystem"
)

type Booking struct {
	BookingID uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"booking_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_user_class" json:"user_id"`
	ClassID   uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_user_class" json:"class_id"`
	Status    BookingStatus  `gorm:"type:booking_status_enum;not null;default:'Confirmed'" json:"status"`
	CreatedAt time.Time      `gorm:"type:timestamptz;not null;default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Class scheduledomain.Class `gorm:"foreignKey:ClassID" json:"-"`
}

type BookingWithClassInfo struct {
	BookingID   uuid.UUID `json:"booking_id"`
	ClassID     uuid.UUID `json:"class_id"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	ServiceName string    `json:"service_name"`
	Instructor  string    `json:"instructor"`
	BranchName  string    `json:"branch_name"`
}

func (Booking) TableName() string {
	return "bookings"
}

type AttendanceStatus string

const (
	AttendanceStatusAttended  AttendanceStatus = "Attended"
	AttendanceStatusNoShow    AttendanceStatus = "NoShow"
	AttendanceStatusCancelled AttendanceStatus = "Cancelled"
)

type AttendanceLog struct {
	AttendanceID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"attendance_id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ServiceID    uuid.UUID `gorm:"type:uuid;not null" json:"service_id"`

	ClassID *uuid.UUID `gorm:"column:schedule_id;type:uuid;index" json:"class_id,omitempty"`

	CheckInTime time.Time        `gorm:"not null;default:now();index" json:"check_in_time"`
	Status      AttendanceStatus `gorm:"type:attendance_status_enum;not null;default:'Attended'" json:"status"`

	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`

	User    authdomain.Users       `gorm:"foreignKey:UserID" json:"-"`
	Service productsdomain.Service `gorm:"foreignKey:ServiceID" json:"-"`
	Class   *scheduledomain.Class  `gorm:"foreignKey:ClassID" json:"-"`
}

func (AttendanceLog) TableName() string {
	return "attendance_log"
}
