package instructordomain

import (
	"context"

	"github.com/google/uuid"
)

type InstructorServiceInterface interface {
	CreateInstructor(ctx context.Context, instructor *Instructor, token string) error
	GetInstructors(ctx context.Context) ([]*Instructor, error)
	DeleteInstructor(ctx context.Context, instructorID uuid.UUID) error
	GetInstructorByID(ctx context.Context, instructorID uuid.UUID) (*Instructor, error)
	UpdateInstructor(ctx context.Context, instructor *Instructor) error

	AssignInstructorsToBranch(ctx context.Context, branchID uuid.UUID, instructorID []uuid.UUID) error
	ListBranchInstructors(ctx context.Context, branchID uuid.UUID) ([]*Instructor, error)
	RemoveInstructorFromBranch(ctx context.Context, branchID uuid.UUID, instructorID uuid.UUID) error
}
