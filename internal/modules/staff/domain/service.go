package staffdomain

import (
	"context"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"github.com/vitalfit/api/pkg/pagination"
)

type StaffServiceInterface interface {
	AssignStaffToBranch(ctx context.Context, branchID uuid.UUID, staffID []uuid.UUID) error
	ListBranchStaffByRole(ctx context.Context, branchID uuid.UUID, roleID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]authdomain.Users, error)
	RemoveStaffFromBranch(ctx context.Context, branchID uuid.UUID, staffID uuid.UUID) error
}
