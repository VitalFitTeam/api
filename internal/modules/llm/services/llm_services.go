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

	userMsg := &llmdomain.Message{
		ConversationID: convo.ConversationID,
		SenderRole:     llmdomain.RoleUser,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	if err := s.store.LLM.SaveMessage(ctx, userMsg); err != nil {
		return "", err
	}

	history, err := s.store.LLM.GetConversationHistory(ctx, convo.ConversationID, 20)
	if err != nil {
		return "", err
	}

	var openaiMsgs []openai.ChatCompletionMessage

	openaiMsgs = append(openaiMsgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "Eres VitalBot, un asistente de gimnasio útil y motivador. Responde de forma concisa.",
	})

	for _, msg := range history {
		role := openai.ChatMessageRoleUser
		if msg.SenderRole == llmdomain.RoleAssistant {
			role = openai.ChatMessageRoleAssistant
		}

		openaiMsgs = append(openaiMsgs, openai.ChatCompletionMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: openaiMsgs,
			// Tools: tools,
		},
	)
	if err != nil {
		s.logger.Errorw("Error calling OpenAI", "error", err)
		return "Lo siento, estoy teniendo problemas de conexión. Intenta más tarde.", nil
	}

	botContent := resp.Choices[0].Message.Content

	botMsg := &llmdomain.Message{
		ConversationID: convo.ConversationID,
		SenderRole:     llmdomain.RoleAssistant,
		Content:        botContent,
		Metadata:       datatypes.JSON([]byte(`{}`)),
		CreatedAt:      time.Now(),
	}

	if err := s.store.LLM.SaveMessage(ctx, botMsg); err != nil {
		s.logger.Error("Error saving bot message", "error", err)
	}

	return botContent, nil
}

func (s *LLMService) GetChatHistory(ctx context.Context, userID uuid.UUID) ([]llmdomain.Message, error) {
	convo, err := s.store.LLM.GetActiveConversation(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error searching for active conversation: %v", err)
	}

	if convo == nil {
		return []llmdomain.Message{}, nil
	}

	history, err := s.store.LLM.GetConversationHistory(ctx, convo.ConversationID, 50)
	if err != nil {
		return nil, fmt.Errorf("error getting history: %v", err)
	}

	return history, nil
}

func (s *LLMService) ResetConversation(ctx context.Context, userID uuid.UUID) error {
	convo, err := s.store.LLM.GetActiveConversation(ctx, userID)
	if err != nil {
		return err
	}

	if convo == nil {
		return nil
	}

	if err := s.store.LLM.DeactivateConversation(ctx, convo.ConversationID); err != nil {
		return fmt.Errorf("error deactivating conversation: %v", err)
	}

	return nil
}
