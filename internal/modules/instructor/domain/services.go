package instructordomain

import (
	"context"

	"github.com/google/uuid"
)

type InstructorServiceInterface interface {
	CreateInstructor(ctx context.Context, instructor *Instructor) error
	GetInstructors(ctx context.Context) ([]*Instructor, error)
	DeleteInstructor(ctx context.Context, instructorID uuid.UUID) error
	GetInstructorByID(ctx context.Context, instructorID uuid.UUID) (*Instructor, error)
	UpdateInstructor(ctx context.Context, instructor *Instructor) error
}
