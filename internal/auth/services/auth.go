package authservices

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// rollbacks user creations if transaction fails
func (h *AuthService) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := h.store.Users.Delete(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) Activate(ctx context.Context, code string) error {
	if err := h.store.Users.Activate(ctx, code); err != nil {
		return err
	}
	return nil

}

func (h *AuthService) GetByEmail(ctx context.Context, email string) (*authdomain.Users, error) {
	users, err := h.store.Users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h *AuthService) CreatePasswordResetToken(ctx context.Context, email string, key string) error {
	user, err := h.store.Users.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if err := h.store.Users.CreatePasswordResetToken(ctx, user.UserID, key, h.store.Config.Mail.Exp); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) DeleteResetToken(ctx context.Context, userID uuid.UUID) error {
	err := h.store.Users.Delete(ctx, userID)
	if err != nil {
		return err
	}
	return nil

}

func (h *AuthService) GenerateToken(user *authdomain.Users) (string, error) {
	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(h.store.Config.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.store.Config.Auth.Token.Iss,
		"aud": h.store.Config.Auth.Token.Iss,
	}
	token, err := h.store.Auth.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *AuthService) ValidateToken(token string) (*jwt.Token, error) {
	return h.store.Auth.ValidateToken(token)
}
