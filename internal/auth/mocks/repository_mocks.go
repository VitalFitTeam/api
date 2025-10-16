package authmocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/auth/domain"
	"gorm.io/gorm"
)

type UserStoreMock struct{}

type RoleStoreMock struct{}

func NewMockUserStore() *UserStoreMock {
	return &UserStoreMock{}
}

func NewMockRoleStore() *RoleStoreMock {
	return &RoleStoreMock{}
}

//USER MOCK FUNCTIONS

func (m *UserStoreMock) Create(ctx context.Context, tx *gorm.DB, user *authdomain.Users) error {
	return nil
}
func (m *UserStoreMock) GetByID(ctx context.Context, userID uuid.UUID) (*authdomain.Users, error) {
	return nil, nil
}
func (m *UserStoreMock) CreateAndInvitate(ctx context.Context, user *authdomain.Users, token string, invitationExp time.Duration) error {
	return nil
}
func (m *UserStoreMock) Delete(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (m *UserStoreMock) Activate(ctx context.Context, code string) error {
	return nil
}
func (m *UserStoreMock) GetByEmail(ctx context.Context, email string) (*authdomain.Users, error) {
	return nil, nil
}
func (m *UserStoreMock) Update(ctx context.Context, user *authdomain.Users) error {
	return nil
}
func (m *UserStoreMock) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, key string, tokenExp time.Duration) error {
	return nil
}
func (m *UserStoreMock) DeleteResetToken(ctx context.Context, userID uuid.UUID) error {
	return nil
}
func (m *UserStoreMock) ResetUserPassword(ctx context.Context, key string, user *authdomain.Users) error {
	return nil
}

// ROLE MOCK FUNCTIONS
func (m *RoleStoreMock) GetByName(ctx context.Context, name string) (*authdomain.Roles, error) {
	return nil, nil
}
