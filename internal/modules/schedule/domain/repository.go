package scheduledomain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ScheduleRepository define las operaciones de acceso a datos para las clases programadas.
type ScheduleRepository interface {

	// CreateClass crea una nueva clase programada.
	CreateClass(ctx context.Context, class *Class) error

	// GetClassesByBranch obtiene todas las clases programadas de una sucursal.
	GetClassesByBranch(ctx context.Context, branchID uuid.UUID) ([]Class, error)

	// GetClassByID obtiene una clase programada específica por su ID.
	GetClassByID(ctx context.Context, classID uuid.UUID) (*Class, error)

	// UpdateClass actualiza los datos de una clase programada.
	UpdateClass(ctx context.Context, class *Class) error

	// DeleteClass elimina (o realiza soft-delete) una clase programada.
	DeleteClass(ctx context.Context, classID uuid.UUID) error

	// GetUpcomingClassesByBranch obtiene las clases futuras (o en curso) de una sucursal.
	GetUpcomingClassesByBranch(ctx context.Context, branchID uuid.UUID) ([]Class, error)

	GetAvailableClassesForBranch(ctx context.Context, branchID uuid.UUID, startTime, endTime time.Time) ([]Class, error)
}
