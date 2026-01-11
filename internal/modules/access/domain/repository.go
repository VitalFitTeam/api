package accessdomain

import (
	"context"

	"github.com/google/uuid"
)

type AccessRepository interface {
	LogAttendance(ctx context.Context, attendance *AttendanceLog) error
	UpdateLogAttendance(ctx context.Context, attendance *AttendanceLog) error
	GetClassAttendanceHistory(ctx context.Context, classID uuid.UUID, startDate, endDate, status *string) ([]*AttendanceLog, error)
}
