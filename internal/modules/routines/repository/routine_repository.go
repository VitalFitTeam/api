package routinerepository

import (
	"context"

	"github.com/google/uuid"
	routinedomain "github.com/vitalfit/api/internal/modules/routines/domain"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
)

type RoutineStore struct {
	db *gorm.DB
}

func NewRoutineStore(db *gorm.DB) *RoutineStore {
	return &RoutineStore{
		db: db,
	}

}

func (s *RoutineStore) CreateRoutine(ctx context.Context, routine *routinedomain.Routine) error {
	return s.db.WithContext(ctx).Create(routine).Error
}

func (s *RoutineStore) AssignRoutine(ctx context.Context, userRoutine *routinedomain.UserRoutine) error {
	return s.db.WithContext(ctx).Create(userRoutine).Error
}

func (s *RoutineStore) GetClientRoutines(ctx context.Context, clientID uuid.UUID) ([]*routinedomain.UserRoutine, error) {
	var routines []*routinedomain.UserRoutine
	err := s.db.WithContext(ctx).
		Preload("Routine.Creator").
		Preload("Routine.RoutineExercises.Exercise").
		Preload("Instructor.User").
		Where("client_id = ? AND status = ?", clientID, routinedomain.StatusActive).
		Order("assigned_date DESC").
		Find(&routines).Error
	return routines, err
}

func (s *RoutineStore) CreateExercise(ctx context.Context, exercise *routinedomain.Exercise) error {
	return s.db.WithContext(ctx).Create(exercise).Error
}

func (s *RoutineStore) GetExercises(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*routinedomain.Exercise, int64, error) {
	var exercises []*routinedomain.Exercise
	var total int64

	query := s.db.WithContext(ctx).Model(&routinedomain.Exercise{})

	if fq.Search != "" {
		query = query.Where("name ILIKE ? OR muscle_group::text ILIKE ?", "%"+fq.Search+"%", "%"+fq.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(fq.Limit).Offset(fq.Page*fq.Limit - fq.Limit).Order("created_at " + fq.Sort).Find(&exercises).Error
	return exercises, total, err
}
