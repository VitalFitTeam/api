package routinedomain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
)

type RoutineServiceInterface interface {
	CreateRoutine(ctx context.Context, routine *Routine) error
	AssignRoutine(ctx context.Context, instructorID, clientID, routineID uuid.UUID, dueDate *time.Time) error
	GetClientRoutines(ctx context.Context, clientID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]*UserRoutine, int64, error)
	CreateExercise(ctx context.Context, exercise *Exercise) error
	GetExercises(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*Exercise, int64, error)
	GetRoutinesByCreator(ctx context.Context, creatorID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]*Routine, int64, error)
	GetRoutineByID(ctx context.Context, routineID uuid.UUID) (*Routine, error)
}
