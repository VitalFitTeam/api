package authdomain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *Users) error
	GetByID(ctx context.Context, userID uuid.UUID) (*Users, error)
	CreateAndInvitate(ctx context.Context, user *Users, token string, invitationExp time.Duration) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Activate(ctx context.Context, code string) error
	GetByEmail(ctx context.Context, email string) (*Users, error)
	Update(ctx context.Context, user *Users) error
	CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, key string, tokenExp time.Duration) error
	DeleteResetToken(ctx context.Context, userID uuid.UUID) error
}

type RolesRepository interface {
	GetByName(ctx context.Context, name string) (*Roles, error)
}
