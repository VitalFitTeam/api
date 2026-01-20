package llmdomain

import (
	"context"

	"github.com/google/uuid"
)

type LLMServiceInterface interface {
	ProcessUserMessage(ctx context.Context, userID uuid.UUID, content string) (string, error)
	GetChatHistory(ctx context.Context, userID uuid.UUID) ([]Message, error)
	ResetConversation(ctx context.Context, userID uuid.UUID) error
}
