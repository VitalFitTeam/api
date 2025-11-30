package accessrepository

import (
	"context"

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
