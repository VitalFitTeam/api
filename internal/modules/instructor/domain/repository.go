package instructordomain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
)

type InstructorRepository interface {
	Create(context.Context, *gorm.DB, *Instructor) error
	GetInstructors(context.Context, pagination.PaginatedFeedQuery) ([]*Instructor, error)
	Delete(context.Context, uuid.UUID) error
	GetByID(context.Context, uuid.UUID) (*Instructor, error)
	Update(context.Context, *Instructor) error
	CreateAndInvitate(ctx context.Context, instructor *Instructor, token string, invitationExp time.Duration) error

	AssignInstructorsToBranch(ctx context.Context, branchID uuid.UUID, instructorID []uuid.UUID) error
	ListBranchInstructors(ctx context.Context, branchID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]*Instructor, error)
	RemoveInstructorFromBranch(ctx context.Context, branchID uuid.UUID, instructorID uuid.UUID) error

	AssignInstructorSpecialty(ctx context.Context, instructorID uuid.UUID, specialties []uuid.UUID) error
	DeleteInstructorSpecialty(ctx context.Context, instructorID uuid.UUID, specialtyID uuid.UUID) error
	AssignInstructorSpecialtyTx(ctx context.Context, tx *gorm.DB, instructorID uuid.UUID, specialties []uuid.UUID) error
}
