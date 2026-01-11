package accessrepository

import (
	"context"

	"github.com/google/uuid"
	accessdomain "github.com/vitalfit/api/internal/modules/access/domain"
	"gorm.io/gorm"
)

type AccessStore struct {
	db *gorm.DB
}

func NewAccessStore(db *gorm.DB) *AccessStore {
	return &AccessStore{
		db: db,
	}
}

func (r *AccessStore) LogAttendance(ctx context.Context, attendance *accessdomain.AttendanceLog) error {
	return r.db.WithContext(ctx).Create(attendance).Error
}

func (r *AccessStore) UpdateLogAttendance(ctx context.Context, attendance *accessdomain.AttendanceLog) error {
	return r.db.WithContext(ctx).Save(attendance).Error
}

func (r *AccessStore) GetClassAttendanceHistory(ctx context.Context, classID uuid.UUID, startDate, endDate, status *string) ([]*accessdomain.AttendanceLog, error) {
	var attendances []*accessdomain.AttendanceLog

	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Service").
		Where("schedule_id = ?", classID)

	// Apply date filters
	if startDate != nil {
		query = query.Where("check_in_time >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("check_in_time <= ?", endDate)
	}

	// Apply status filter
	if status != nil {
		query = query.Where("status = ?", status)
	}

	// Order by check-in time descending (most recent first)
	query = query.Order("check_in_time DESC")

	err := query.Find(&attendances).Error
	if err != nil {
		return nil, err
	}

	return attendances, nil
}
