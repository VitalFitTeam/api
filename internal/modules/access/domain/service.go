package accessdomain

import (
	"context"

	"github.com/google/uuid"
)

type AcessServiceInterface interface {
	ProcessCheckIn(ctx context.Context, userID, branchID uuid.UUID) (*CheckInResponse, error)
}
