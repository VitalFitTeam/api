package instructordomain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstructorRepository interface {
	Create(context.Context, *gorm.DB, *Instructor) error
	GetInstructors(context.Context) ([]*Instructor, error)
	Delete(context.Context, uuid.UUID) error
	GetByID(context.Context, uuid.UUID) (*Instructor, error)
	Update(context.Context, *Instructor) error

	CreateAndInvitate(ctx context.Context, instructor *Instructor, token string, invitationExp time.Duration) error
}
