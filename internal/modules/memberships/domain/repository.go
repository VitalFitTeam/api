package membershipsdomain

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipsRepository interface {
	CreateMembershipType(ctx context.Context, membership *MembershipType) error
	CreateMembershipTypeTX(ctx context.Context, tx *gorm.DB, membership *MembershipType) error
	UpdateMembershipType(ctx context.Context, membership *MembershipType) error
	DeleteMembershipType(ctx context.Context, id uuid.UUID) error
	GetMembershipTypeByID(ctx context.Context, id uuid.UUID) (*MembershipType, error)
	GetMembershipTypes(ctx context.Context) ([]*MembershipType, error)
}
