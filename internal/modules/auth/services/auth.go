package authservices

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vitalfit/api/config"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/mailer"
)

type AuthService struct {
	store  store.Storage
	config config.Config
	auth   authdomain.Authenticator
	Mailer mailer.Client
}

func NewAuthServices(store store.Storage, cfg config.Config, auth authdomain.Authenticator, mailer mailer.Client) *AuthService {
	return &AuthService{
		store:  store,
		config: cfg,
		auth:   auth,
		Mailer: mailer,
	}
}

func (s *AuthService) RegisterUserClient(ctx context.Context, user *authdomain.Users, token string) error {
	role, error := s.store.Roles.GetByName(ctx, "client")
	client_profile := &authdomain.ClientProfiles{
		UserID:   user.UserID,
		Category: authdomain.ClientCategoryNew,
	}
	if error != nil {
		return error
	}
	user.RoleID = role.RoleID
	user.ClientProfile = *client_profile
	if err := s.store.Users.CreateAndInvitate(ctx, user, token, s.config.Mail.Exp); err != nil {
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
	if err := s.store.Users.CreateAndInvitate(ctx, user, token, s.config.Mail.Exp); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) MailSender(ctx context.Context, user *authdomain.Users, key string, template string) (int, error) {

	//mail -> fail -> roll back -> create invite

	isProdEnv := h.config.Env == "production"
	vars := struct {
		Username string
		CODE     string
	}{
		Username: user.FirstName,
		CODE:     key,
	}

	// send mail
	status, err := h.Mailer.Send(template, user.FirstName, user.Email, vars, !isProdEnv)
	if err != nil {
		return status, err
	}

	return status, err
}

func (h *AuthService) MailSenderStaff(ctx context.Context, user *authdomain.Users, token string, template string) (int, error) {

	//mail -> fail -> roll back -> create invite

	isProdEnv := h.config.Env == "production"
	vars := struct {
		Username       string
		ActivationLink string
	}{
		Username:       user.FirstName,
		ActivationLink: fmt.Sprintf("%s/activate/%s", h.config.FrontURL, token),
	}

	// send mail
	status, err := h.Mailer.Send(template, user.FirstName, user.Email, vars, !isProdEnv)
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

func (h *AuthService) ActivateStaff(ctx context.Context, token string, password string) error {
	if err := h.store.Users.ActivateUserStaff(ctx, token, password); err != nil {
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
	if err := h.store.Users.CreatePasswordResetToken(ctx, user.UserID, key, h.config.Mail.Exp); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) DeleteResetToken(ctx context.Context, userID uuid.UUID) error {
	err := h.store.Users.DeleteResetToken(ctx, userID)
	if err != nil {
		return err
	}
	return nil

}

func (h *AuthService) GenerateToken(user *authdomain.Users) (string, error) {
	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(h.config.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.Token.Iss,
		"aud": h.config.Auth.Token.Aud,
	}
	token, err := h.auth.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (h *AuthService) ValidateToken(token string) (*jwt.Token, error) {
	return h.auth.ValidateToken(token)
}

func (h *AuthService) ResetPassword(ctx context.Context, key string, user *authdomain.Users) error {
	if err := h.store.Users.ResetUserPassword(ctx, key, user); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) ValidateResetToken(ctx context.Context, key string) error {
	if err := h.store.Users.ValidateResetToken(ctx, key); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) UpgradePassword(ctx context.Context, user *authdomain.Users) error {
	if err := h.store.Users.UpgradePassword(ctx, user); err != nil {
		return err
	}
	return nil
}

func (h *AuthService) GenerateQrJwtToken(ctx context.Context, user *authdomain.Users) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(30 * time.Second).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.Token.Iss,
		"aud": h.config.Auth.Token.Aud,
	}
	token, err := h.auth.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
