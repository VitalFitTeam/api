package instructordomain

import (
	"context"

	"github.com/google/uuid"
)

type InstructorRepository interface {
	Create(context.Context, *Instructor) error
	GetInstructors(context.Context) ([]*Instructor, error)
	Delete(context.Context, uuid.UUID) error
	GetByID(context.Context, uuid.UUID) (*Instructor, error)
	Update(context.Context, *Instructor) error
}
