package authservices

import (
	"context"

	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"github.com/vitalfit/api/internal/store"
)

type AuthService struct {
	store store.Storage
}

func NewAuthServices(store store.Storage) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, user authdomain.Users, roleName string) error {
	role, error := s.store.Roles.GetByName(ctx, roleName)
	if error != nil {
		return error
	}
	user.RoleID = role.RoleID
	if err := s.store.Users.Create(ctx, user); err != nil {
		return err
	}

	return nil
}
