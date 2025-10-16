package store

import (
	authmocks "github.com/vitalfit/api/internal/auth/mocks"
)

func NewMockStore() Storage {
	return Storage{
		Users: authmocks.NewMockUserStore(),
		Roles: authmocks.NewMockRoleStore(),
	}
}
