package staffservice

import (
	"context"

	"github.com/google/uuid"
	authdomain "github.com/vitalfit/api/internal/modules/auth/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/pagination"
)

type StaffService struct {
	store store.Storage
}

func NewStaffService(store store.Storage) *StaffService {
	return &StaffService{store: store}
}

func (s *StaffService) AssignStaffToBranch(ctx context.Context, branchID uuid.UUID, staffID []uuid.UUID) error {
	return s.store.Staff.AssignStaffToBranch(ctx, branchID, staffID)
}

func (s *StaffService) ListBranchStaffByRole(ctx context.Context, branchID uuid.UUID, roleID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]authdomain.Users, error) {
	return s.store.Staff.ListBranchStaffByRole(ctx, branchID, roleID, fq)
}

func (s *StaffService) RemoveStaffFromBranch(ctx context.Context, branchID uuid.UUID, staffID uuid.UUID) error {
	return s.store.Staff.RemoveStaffFromBranch(ctx, branchID, staffID)
}
