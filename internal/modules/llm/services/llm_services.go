package llmservices

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	llmdomain "github.com/vitalfit/api/internal/modules/llm/domain"
	"github.com/vitalfit/api/internal/store"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type LLMService struct {
	client *openai.Client
	store  store.Storage
	logger *zap.SugaredLogger
}

func NewLLMService(client *openai.Client, store store.Storage, logger *zap.SugaredLogger) *LLMService {
	return &LLMService{
		client: client,
		store:  store,
		logger: logger,
	}
}

func (s *LLMService) ProcessUserMessage(ctx context.Context, userID uuid.UUID, content string) (string, error) {
	convo, err := s.store.LLM.GetActiveConversation(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("error searching for conversation: %v", err)
	}
	if convo == nil {
		convo, err = s.store.LLM.CreateConversation(ctx, userID)
		if err != nil {
			return "", fmt.Errorf("error creating conversation: %v", err)
		}
	}

	// 2. Guardar mensaje del Usuario en BD (Antes de llamar a la IA)
	userMsg := &llmdomain.Message{
		ConversationID: convo.ConversationID,
		SenderRole:     llmdomain.RoleUser,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	if err := s.store.LLM.SaveMessage(ctx, userMsg); err != nil {
		return "", err
	}

	// 3. Construir Historial para OpenAI
	// Recuperamos mensajes previos de la BD
	history, err := s.store.LLM.GetConversationHistory(ctx, convo.ConversationID, 20)
	if err != nil {
		return "", err
	}

	var openaiMsgs []openai.ChatCompletionMessage

	// System Prompt (Personalidad)
	openaiMsgs = append(openaiMsgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "Eres VitalBot, un asistente de gimnasio útil y motivador. Responde de forma concisa.",
	})

	// Agregar historial de BD
	for _, msg := range history {
		role := openai.ChatMessageRoleUser
		if msg.SenderRole == llmdomain.RoleAssistant {
			role = openai.ChatMessageRoleAssistant
		}
		// (Aquí podrías manejar role 'tool' si usas tools)

		openaiMsgs = append(openaiMsgs, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// 4. Llamar a OpenAI
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: openaiMsgs,
			// Tools: tools, (Aquí inyectarías tus tools de rutinas)
		},
	)
	if err != nil {
		s.logger.Errorw("Error calling OpenAI", "error", err)
		return "Lo siento, estoy teniendo problemas de conexión. Intenta más tarde.", nil
	}

	botContent := resp.Choices[0].Message.Content

	// 5. Guardar Respuesta del Bot en BD
	// (Opcional: guardar tokens usados en metadata)
	botMsg := &llmdomain.Message{
		ConversationID: convo.ConversationID,
		SenderRole:     llmdomain.RoleAssistant,
		Content:        botContent,
		Metadata:       datatypes.JSON([]byte(`{}`)), // Aquí podrías poner tokens
		CreatedAt:      time.Now(),
	}

	if err := s.store.LLM.SaveMessage(ctx, botMsg); err != nil {
		s.logger.Error("Error saving bot message", "error", err)
		// No retornamos error aquí porque el usuario ya leyó la respuesta
	}

	return botContent, nil
}

func (s *LLMService) GetChatHistory(ctx context.Context, userID uuid.UUID) ([]llmdomain.Message, error) {
	// 1. Buscar conversación activa
	convo, err := s.store.LLM.GetActiveConversation(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error buscando conversación activa: %v", err)
	}

	// Si no hay conversación activa, retornamos lista vacía (no es error)
	if convo == nil {
		return []llmdomain.Message{}, nil
	}

	// 2. Traer mensajes (Limitamos a 50 para no sobrecargar la vista inicial)
	// Asumimos que GetConversationHistory devuelve los mensajes ordenados cronológicamente
	history, err := s.store.LLM.GetConversationHistory(ctx, convo.ConversationID, 50)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo historial: %v", err)
	}

	return history, nil
}

// ResetConversation marca la conversación actual como inactiva (Soft Delete lógico)
func (s *LLMService) ResetConversation(ctx context.Context, userID uuid.UUID) error {
	// 1. Buscar conversación activa
	convo, err := s.store.LLM.GetActiveConversation(ctx, userID)
	if err != nil {
		return err
	}

	// Si no existe, no hay nada que borrar
	if convo == nil {
		return nil
	}

	// 2. Desactivar en BD (is_active = false)
	// Necesitas asegurarte de que este método exista en tu Repo
	if err := s.store.LLM.DeactivateConversation(ctx, convo.ConversationID); err != nil {
		return fmt.Errorf("error desactivando conversación: %v", err)
	}

	return nil
}
