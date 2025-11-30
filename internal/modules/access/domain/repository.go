package accessdomain

import "context"

type AccessRepository interface {
	LogAttendance(ctx context.Context, attendance *AttendanceLog) error
	UpdateLogAttendance(ctx context.Context, attendance *AttendanceLog) error
}
