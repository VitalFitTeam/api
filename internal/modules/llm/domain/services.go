package llmdomain

import (
	"context"

	"github.com/google/uuid"
)

type LLMServiceInterface interface {
	ProcessUserMessage(ctx context.Context, userID uuid.UUID, content string) (string, error)
}
