package staffdomain

import (
	"context"

	"github.com/google/uuid"
)

type StaffRepository interface {
	AssignStaffToBranch(ctx context.Context, branchID uuid.UUID, staffID []uuid.UUID) error
}
