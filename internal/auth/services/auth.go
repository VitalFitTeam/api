package authservices

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vitalfit/api/internal/store"
)

type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type AuthServicesInterface interface {
	RegisterUser(ctx context.Context)
}

type AuthService struct {
	store store.Storage
}

func NewAuthServices(store store.Storage) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context) {

}
