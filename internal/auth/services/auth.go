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

func (s *AuthService) RegisterUserClient(ctx context.Context, user authdomain.Users, token string) error {
	role, error := s.store.Roles.GetByName(ctx, "client")
	if error != nil {
		return error
	}
	user.RoleID = role.RoleID
	if err := s.store.Users.CreateAndInvitate(ctx, user, token, s.store.Config.Mail.Exp); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) RegisterUserStaff(ctx context.Context, user authdomain.Users, token string, roleName string) error {
	role, error := s.store.Roles.GetByName(ctx, "client")
	if error != nil {
		return error
	}
	user.RoleID = role.RoleID
	if err := s.store.Users.CreateAndInvitate(ctx, user, token, s.store.Config.Mail.Exp); err != nil {
		return err
	}
	return nil
}
