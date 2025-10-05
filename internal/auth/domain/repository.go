package authdomain

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource already exists")
)

type UserRepository interface {
	Create(context.Context, Users) error
	GetUser() error
}

type RolesRepository interface {
	GetByName(ctx context.Context, name string) (*Roles, error)
}
