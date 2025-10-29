package instructorrepository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/pkg/db"
	"gorm.io/gorm"
)

type InstructorStore struct {
	db *gorm.DB
}

func NewInstructorStore(db *gorm.DB) *InstructorStore {
	return &InstructorStore{db: db}
}

func (s *InstructorStore) Create(ctx context.Context, instructor *instructordomain.Instructor) error {
	err := s.db.WithContext(ctx).Create(instructor).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique constraint violation
			return shared_errors.ErrConflict
		}
		return err
	}
	return nil
}

func (s *InstructorStore) GetInstructors(ctx context.Context) ([]*instructordomain.Instructor, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	var instructors []*instructordomain.Instructor
	err := s.db.WithContext(ctx).Find(&instructors).Error
	if err != nil {
		return nil, err
	}
	return instructors, nil
}

func (s *InstructorStore) Delete(ctx context.Context, instructorID uuid.UUID) error {
	err := s.db.WithContext(ctx).Delete(&instructordomain.Instructor{}, instructorID).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return shared_errors.ErrNotFound
		default:
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return shared_errors.ErrConflict
			}
			return err
		}
	}
	return nil
}

func (s *InstructorStore) GetByID(ctx context.Context, instructorID uuid.UUID) (*instructordomain.Instructor, error) {
	ctx, cancel := context.WithTimeout(ctx, db.QueryTimeoutDuration)
	defer cancel()
	var instructor *instructordomain.Instructor
	err := s.db.WithContext(ctx).Where("instructor_id = ?", instructorID).First(&instructor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return instructor, nil
}

func (s *InstructorStore) Update(ctx context.Context, instructor *instructordomain.Instructor) error {
	result := s.db.WithContext(ctx).Model(&instructordomain.Instructor{}).Where("instructor_id = ?", instructor.InstructorID).Updates(instructor)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}
	return nil
}
