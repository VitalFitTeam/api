package membershipsservice

import (
	"context"

	"github.com/google/uuid"
	membershipsdomain "github.com/vitalfit/api/internal/modules/memberships/domain"
	"github.com/vitalfit/api/internal/store"
)

// MembershipService implementa la lógica de negocio para los tipos de membresía.
type MembershipService struct {
	store store.Storage
}

// NewMembershipService crea una nueva instancia del servicio de membresías.
func NewMembershipService(store store.Storage) *MembershipService {
	return &MembershipService{
		store: store,
	}
}

// CreateMembershipType crea un nuevo tipo de membresía.
func (s *MembershipService) CreateMembershipType(ctx context.Context, membership *membershipsdomain.MembershipType) error {
	err := s.store.Membership.CreateMembershipType(ctx, membership)
	if err != nil {
		return err
	}
	return nil
}

// UpdateMembershipType actualiza un tipo de membresía existente.
func (s *MembershipService) UpdateMembershipType(ctx context.Context, membership *membershipsdomain.MembershipType) error {
	err := s.store.Membership.UpdateMembershipType(ctx, membership)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMembershipType realiza una eliminación lógica (soft delete) del tipo de membresía.
func (s *MembershipService) DeleteMembershipType(ctx context.Context, id uuid.UUID) error {
	err := s.store.Membership.DeleteMembershipType(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// GetMembershipTypeByID obtiene un tipo de membresía por su ID.
func (s *MembershipService) GetMembershipTypeByID(ctx context.Context, id uuid.UUID) (*membershipsdomain.MembershipType, error) {
	membership, err := s.store.Membership.GetMembershipTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return membership, nil
}

// GetMembershipTypes devuelve todos los tipos de membresía activos.
func (s *MembershipService) GetMembershipTypes(ctx context.Context) ([]*membershipsdomain.MembershipType, error) {
	memberships, err := s.store.Membership.GetMembershipTypes(ctx)
	if err != nil {
		return nil, err
	}
	return memberships, nil
}
