package authdomain

import (
	"context"
)

type UserRepository interface {
	Create(context.Context, Users) error
	GetUser() error
}

type RolesRepository interface {
	GetByName(ctx context.Context, name string) (*Roles, error)
}
