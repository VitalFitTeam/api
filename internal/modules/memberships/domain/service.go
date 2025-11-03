package membershipsdomain

import (
	"context"

	"github.com/google/uuid"
)

type MembershipsServiceInterface interface {
	CreateMembershipType(ctx context.Context, membership *MembershipType) error
	UpdateMembershipType(ctx context.Context, membership *MembershipType) error
	DeleteMembershipType(ctx context.Context, id uuid.UUID) error
	GetMembershipTypeByID(ctx context.Context, id uuid.UUID) (*MembershipType, error)
	GetMembershipTypes(ctx context.Context) ([]*MembershipType, error)
}
