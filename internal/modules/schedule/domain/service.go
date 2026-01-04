package scheduledomain

import (
	"context"

	"github.com/google/uuid"
)

type ScheduleServiceInterface interface {
	CreateClass(ctx context.Context, class *Class) error
	GetClassesByBranch(ctx context.Context, branchID uuid.UUID) ([]Class, error)
	GetClassByID(ctx context.Context, classID uuid.UUID) (*Class, error)
	UpdateClass(ctx context.Context, class *Class) error
	DeleteClass(ctx context.Context, classID uuid.UUID) error
	GetUpcomingClassesByBranch(ctx context.Context, branchID uuid.UUID) ([]Class, error)
}
