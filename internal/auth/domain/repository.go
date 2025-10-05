package authdomain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user Users) error
	GetUser() error
	CreateAndInvitate(ctx context.Context, user Users, token string, invitationExp time.Duration) error
}

type RolesRepository interface {
	GetByName(ctx context.Context, name string) (*Roles, error)
}
