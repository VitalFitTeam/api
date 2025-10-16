package authdomain

import (
	"context"

	"github.com/gin-gonic/gin"
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
	MailSender(ctx context.Context, user *Users, key string, template string) (int, error)
	Activate(ctx context.Context, code string) error
	GenerateToken(user *Users) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	CreatePasswordResetToken(ctx context.Context, email string, key string) error
	DeleteResetToken(context.Context, uuid.UUID) error
	ResetPassword(ctx context.Context, key string, user *Users) error
}

type UserServicesInterface interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*Users, error)
	Update(ctx context.Context, user *Users) error
	GetByEmail(ctx context.Context, email string) (*Users, error)
	GetUserFromContext(c *gin.Context) *Users
	GetRoleByName(ctx context.Context, name string) (*Roles, error)
}
