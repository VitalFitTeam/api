package instructorservices

import (
	"context"

	"github.com/google/uuid"
	instructordomain "github.com/vitalfit/api/internal/modules/instructor/domain"
	shared_errors "github.com/vitalfit/api/internal/shared/errors"
	"github.com/vitalfit/api/internal/store"
)

type InstructorServices struct {
	store store.Storage
}

func NewInstructorServices(store store.Storage) *InstructorServices {
	return &InstructorServices{
		store: store,
	}
}

func (s *InstructorServices) CreateInstructor(ctx context.Context, instructor *instructordomain.Instructor) error {
	role, err := s.store.Roles.GetByName(ctx, "instructor")
	if err != nil {
		return err
	}
	instructor.User.RoleID = role.RoleID
	if err = s.store.Instructor.Create(ctx, instructor); err != nil {
		return err
	}
	return nil
}

func (s *InstructorServices) GetInstructors(ctx context.Context) ([]*instructordomain.Instructor, error) {
	instructors, err := s.store.Instructor.GetInstructors(ctx)
	if err != nil {
		return nil, err
	}
	return instructors, nil
}

func (s *InstructorServices) DeleteInstructor(ctx context.Context, instructorID uuid.UUID) error {
	err := s.store.Instructor.Delete(ctx, instructorID)
	if err != nil {
		if err == shared_errors.ErrNotFound {
			return shared_errors.ErrNotFound
		}
		return err
	}
	return nil
}

func (s *InstructorServices) GetInstructorByID(ctx context.Context, instructorID uuid.UUID) (*instructordomain.Instructor, error) {
	instructor, err := s.store.Instructor.GetByID(ctx, instructorID)
	if err != nil {
		if err == shared_errors.ErrNotFound {
			return nil, shared_errors.ErrNotFound
		}
		return nil, err
	}
	return instructor, nil
}

func (s *InstructorServices) UpdateInstructor(ctx context.Context, instructor *instructordomain.Instructor) error {
	err := s.store.Instructor.Update(ctx, instructor)
	if err != nil {
		if err == shared_errors.ErrNotFound {
			return shared_errors.ErrNotFound
		}
		return err
	}
	return nil
}
