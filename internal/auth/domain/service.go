package authdomain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type AuthServicesInterface interface {
	RegisterUserClient(ctx context.Context, user *Users, token string) error
	RegisterUserStaff(ctx context.Context, user *Users, token string, roleName string) error
	Delete(context.Context, uuid.UUID) error
	MailSender(ctx context.Context, user *Users, key string) (int, error)
	Activate(ctx context.Context, code string) error
}
