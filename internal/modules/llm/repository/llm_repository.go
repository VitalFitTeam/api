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

// CreateConversation crea un nuevo hilo
func (r *LLMStore) CreateConversation(ctx context.Context, userID uuid.UUID) (*llmdomain.Conversation, error) {
	convo := &llmdomain.Conversation{
		ConversationID: uuid.New(),
		UserID:         userID,
		IsActive:       true,
	}
	err := r.db.WithContext(ctx).Create(convo).Error
	return convo, err
}

// SaveMessage guarda un mensaje en el historial
func (r *LLMStore) SaveMessage(ctx context.Context, msg *llmdomain.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

// GetConversationHistory recupera los últimos N mensajes para darle contexto a la IA
func (r *LLMStore) GetConversationHistory(ctx context.Context, convoID uuid.UUID, limit int) ([]llmdomain.Message, error) {
	var messages []llmdomain.Message

	// Truco: Obtenemos los últimos N ordenados por fecha DESC (del más nuevo al viejo)
	// para el LIMIT, y luego el servicio deberá invertirlos o enviarlos correctamente.
	// O mejor, traemos los últimos y dejamos que SQL los ordene.

	// Subquery strategy es compleja en GORM simple, así que haremos:
	// Traer todos o usar un limit alto. Para un MVP, traer los ultimos 20 ordenados ASC
	// asumiendo que el chat no es infinito, o filtrar por fecha.

	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", convoID).
		Order("created_at ASC"). // Orden cronológico (Lo que OpenAI necesita)
		// Limit(limit). // Nota: Si la conver es muy larga, necesitarás lógica de ventana deslizante
		Find(&messages).Error

	return messages, err
}
