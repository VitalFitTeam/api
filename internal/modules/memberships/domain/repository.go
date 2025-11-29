package membershipsdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
	"gorm.io/gorm"
)

type MembershipsRepository interface {
	CreateMembershipType(ctx context.Context, membership *MembershipType) error
	CreateMembershipTypeTX(ctx context.Context, tx *gorm.DB, membership *MembershipType) error
	UpdateMembershipType(ctx context.Context, membership *MembershipType) error
	DeleteMembershipType(ctx context.Context, id uuid.UUID) error
	GetMembershipTypeByID(ctx context.Context, id uuid.UUID) (*MembershipType, error)
	GetMembershipTypes(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*MembershipType, error)
	GetMembershipTypesFTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error)
	GetSummary(ctx context.Context) (*MembershipSummary, error)

	UpdateClientMembership(ctx context.Context, membership *ClientMembership) error
	ClientHasActiveMembership(ctx context.Context, clientID uuid.UUID) (bool, error)
	GetClientMembership(ctx context.Context, clientID uuid.UUID) (*ClientMembership, error)
}
