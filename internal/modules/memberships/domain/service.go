package membershipsdomain

import (
	"context"

	"github.com/google/uuid"
	"github.com/vitalfit/api/pkg/pagination"
)

type MembershipsServiceInterface interface {
	CreateMembershipType(ctx context.Context, membership *MembershipType) error
	UpdateMembershipType(ctx context.Context, membership *MembershipType) error
	DeleteMembershipType(ctx context.Context, id uuid.UUID) error
	GetMembershipTypeByID(ctx context.Context, id uuid.UUID) (*MembershipType, error)
	GetMembershipTypes(ctx context.Context, fq pagination.PaginatedFeedQuery) ([]*MembershipType, error)
	GetMembershipTypesFTotal(ctx context.Context, fq pagination.PaginatedFeedQuery) (int64, error)
	GetSummary(ctx context.Context) (*MembershipSummary, error)
}
