package authservices

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"github.com/vitalfit/api/internal/store"
)

type UserService struct {
	store store.Storage
}

func NewUserService(store store.Storage) *UserService {
	return &UserService{
		store: store,
	}
}

func (h *UserService) GetByEmail(ctx context.Context, email string) (*authdomain.Users, error) {
	users, err := h.store.Users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h *UserService) GetByID(ctx context.Context, userID uuid.UUID) (*authdomain.Users, error) {
	users, err := h.store.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (h *UserService) Update(ctx context.Context, user *authdomain.Users) error {
	if err := h.store.Users.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (h *UserService) GetUserFromContext(c *gin.Context) *authdomain.Users {
	user, ok := c.Value("user").(*authdomain.Users)
	if !ok {
		return nil
	}
	return user
}

func (h *UserService) GetRoleByName(ctx context.Context, name string) (*authdomain.Roles, error) {
	role, err := h.store.Roles.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return role, nil
}
