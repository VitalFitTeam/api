package authservices

import (
	"context"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
)

type AuthService struct {
	store store.Storage
}

func NewAuthServices(store store.Storage) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) RegisterUserClient(ctx context.Context, user *authdomain.Users, token string) error {
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

func (s *AuthService) RegisterUserStaff(ctx context.Context, user *authdomain.Users, token string, roleName string) error {
	role, error := s.store.Roles.GetByName(ctx, roleName)
	if error != nil {
		return error
	}
	user.RoleID = role.RoleID
	if err := s.store.Users.CreateAndInvitate(ctx, user, token, s.store.Config.Mail.Exp); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) MailSender(ctx context.Context, user *authdomain.Users, key string) (int, error) {

	//mail -> fail -> roll back -> create invite

	isProdEnv := h.store.Env == "production"
	vars := struct {
		Username       string
		ActivationCODE string
	}{
		Username:       user.FirstName,
		ActivationCODE: key,
	}

	// send mail
	status, err := h.store.Mailer.Send(mailer.UserWelcomeTemplate, user.FirstName, user.Email, vars, !isProdEnv)
	if err != nil {
		return status, err
	}

	return status, err
}

func (h *AuthService) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := h.store.Users.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}
