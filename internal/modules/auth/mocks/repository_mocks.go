package authmocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type UserStoreMock struct {
	mock.Mock
}

type RoleStoreMock struct {
	mock.Mock
}

func NewMockUserStore() *UserStoreMock {
	return &UserStoreMock{}
}

func NewMockRoleStore() *RoleStoreMock {
	return &RoleStoreMock{}
}

//USER MOCK FUNCTIONS

func (m *UserStoreMock) Create(ctx context.Context, tx *gorm.DB, user *authdomain.Users) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *UserStoreMock) GetByID(ctx context.Context, userID uuid.UUID) (*authdomain.Users, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.Users), args.Error(1)
}

func (m *UserStoreMock) CreateAndInvitate(ctx context.Context, user *authdomain.Users, token string, invitationExp time.Duration) error {
	args := m.Called(ctx, user, token, invitationExp)
	return args.Error(0)
}

func (m *UserStoreMock) Delete(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *UserStoreMock) Activate(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

func (m *UserStoreMock) GetByEmail(ctx context.Context, email string) (*authdomain.Users, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.Users), args.Error(1)
}

func (m *UserStoreMock) Update(ctx context.Context, user *authdomain.Users) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserStoreMock) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, key string, tokenExp time.Duration) error {
	args := m.Called(ctx, userID, key, tokenExp)
	return args.Error(0)
}

func (m *UserStoreMock) DeleteResetToken(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *UserStoreMock) ResetUserPassword(ctx context.Context, key string, user *authdomain.Users) error {
	args := m.Called(ctx, key, user)
	return args.Error(0)
}

// ROLE MOCK FUNCTIONS
func (m *RoleStoreMock) GetByName(ctx context.Context, name string) (*authdomain.Roles, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.Roles), args.Error(1)
}
