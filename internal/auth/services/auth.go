package authservices

import (
	"context"

	"github.com/vitalfit/api/internal/store"
)

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
