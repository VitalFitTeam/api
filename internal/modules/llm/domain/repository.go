package llmdomain

import (
	"context"

	"github.com/google/uuid"
)

type LLMRepository interface {
	GetActiveConversation(ctx context.Context, userID uuid.UUID) (*Conversation, error)
	CreateConversation(ctx context.Context, userID uuid.UUID) (*Conversation, error)
	SaveMessage(ctx context.Context, msg *Message) error
	GetConversationHistory(ctx context.Context, convoID uuid.UUID, limit int) ([]Message, error)
}
