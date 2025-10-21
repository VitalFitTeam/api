package store

import authmocks "github.com/vitalfit/api/internal/modules/auth/mocks"

func NewMockStore() Storage {
	return Storage{
		Users: authmocks.NewMockUserStore(),
		Roles: authmocks.NewMockRoleStore(),
	}
}
