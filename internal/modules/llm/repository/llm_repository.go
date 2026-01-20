package llmrepository

import (
	"context"

	"github.com/google/uuid"
	llmdomain "github.com/vitalfit/api/internal/modules/llm/domain"
	"gorm.io/gorm"
)

type LLMStore struct {
	db *gorm.DB
}

func NewLLMStore(db *gorm.DB) *LLMStore {
	return &LLMStore{db: db}
}

func (r *LLMStore) GetActiveConversation(ctx context.Context, userID uuid.UUID) (*llmdomain.Conversation, error) {
	var convo llmdomain.Conversation
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		First(&convo).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &convo, nil
}

func (r *LLMStore) CreateConversation(ctx context.Context, userID uuid.UUID) (*llmdomain.Conversation, error) {
	convo := &llmdomain.Conversation{
		ConversationID: uuid.New(),
		UserID:         userID,
		IsActive:       true,
	}
	err := r.db.WithContext(ctx).Create(convo).Error
	return convo, err
}

func (r *LLMStore) SaveMessage(ctx context.Context, msg *llmdomain.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *LLMStore) GetConversationHistory(ctx context.Context, convoID uuid.UUID, limit int) ([]llmdomain.Message, error) {
	var messages []llmdomain.Message

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", convoID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}

func (r *LLMStore) DeactivateConversation(ctx context.Context, convoID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&llmdomain.Conversation{}).
		Where("conversation_id = ?", convoID).
		Update("is_active", false).Error
}
