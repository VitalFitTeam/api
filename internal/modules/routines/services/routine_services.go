package routineservices

import (
	"context"
	"time"

	"github.com/google/uuid"
	routinedomain "github.com/vitalfit/api/internal/modules/routines/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/pagination"
)

type RoutineService struct {
	store store.Storage
}

func NewRoutineService(store store.Storage) *RoutineService {
	return &RoutineService{
		store: store,
	}

}

func (s *RoutineService) CreateRoutine(ctx context.Context, routine *routinedomain.Routine) error {
	return s.store.Routine.CreateRoutine(ctx, routine)
}

func (s *RoutineService) AssignRoutine(ctx context.Context, instructorUserID, clientID, routineID uuid.UUID, dueDate *time.Time) error {
	// Resolve InstructorID from UserID
	instructor, err := s.store.Instructor.GetByUserID(ctx, instructorUserID)
	if err != nil {
		return err
	}

	assignment := &routinedomain.UserRoutine{
		ClientID:     clientID,
		InstructorID: instructor.InstructorID,
		RoutineID:    routineID,
		Status:       routinedomain.StatusActive,
		DueDate:      dueDate,
		AssignedDate: time.Now(),
		IsActive:     true,
	}

	return s.store.Routine.AssignRoutine(ctx, assignment)
}

func (s *RoutineService) GetClientRoutines(ctx context.Context, clientID uuid.UUID) ([]*routinedomain.UserRoutine, error) {
	return s.store.Routine.GetClientRoutines(ctx, clientID)
}

func (s *RoutineService) CreateExercise(ctx context.Context, exercise *routinedomain.Exercise) error {
	return s.store.Routine.CreateExercise(ctx, exercise)
}

func (s *RoutineService) GetExercises(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*routinedomain.Exercise, int64, error) {
	return s.store.Routine.GetExercises(ctx, fq)
}
