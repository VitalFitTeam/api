package authdomain

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
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
	MailSenderStaff(ctx context.Context, user *Users, token string, template string) (int, error)
	Activate(ctx context.Context, code string) error
	ActivateStaff(ctx context.Context, token string, password string) error
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
	GetBranchAdmins(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*Users, error)
	GetUsers(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*Users, error)
	GetClients(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*Users, error)
	UpdateClient(ctx context.Context, user *Users) error
	UpdateStaff(ctx context.Context, user *Users, roleName string) error
	//roles
	GetRoleByName(ctx context.Context, name string) (*Roles, error)
	RoleHasPermission(ctx context.Context, roleID uuid.UUID, permission string) (bool, error)
	GetRoles(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*Roles, error)
	GetRolesFTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error)
	CreateRole(ctx context.Context, role *Roles) error
	GetRoleByID(ctx context.Context, roleID uuid.UUID) (*Roles, error)
	UpdateRole(ctx context.Context, role *Roles) error
	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	GetPermissions(ctx context.Context) ([]*Permission, error)
	AssignRolePermission(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	DeleteRolePermission(ctx context.Context, roleID uuid.UUID, permissionID []uuid.UUID) error
}
