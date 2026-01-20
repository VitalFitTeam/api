package llmservices

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	llmdomain "github.com/vitalfit/api/internal/modules/llm/domain"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/pagination"
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

	fq := pagination.PaginatedFeedQuery{Limit: 20, Page: 1, Sort: "desc"}
	history, _, err := s.store.LLM.GetConversationHistory(ctx, convo.ConversationID, fq)
	if err != nil {
		return "", err
	}

	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}

	var openaiMsgs []openai.ChatCompletionMessage

	systemPrompt := fmt.Sprintf(`Eres VitalBot, el asistente virtual de VitalFit.
Tu misión es motivar a los usuarios y ayudarles con la gestión de su gimnasio.
Fecha y hora actual: %s.

Directrices:
1. Responde de forma concisa y amigable.
2. Usa las herramientas disponibles para consultar horarios, gestionar reservas y crear rutinas.
3. Si te piden clases para "hoy" o "mañana", usa la fecha actual como referencia.
4. Si falta información para una herramienta (ej. ID de sucursal), pregunta al usuario o usa 'get_all_branches' para guiarlo.`, time.Now().Format("2006-01-02 15:04"))

	openaiMsgs = append(openaiMsgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemPrompt,
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
			Tools:    tools,
		},
	)
	if err != nil {
		s.logger.Errorw("Error calling OpenAI", "error", err)
		return "Lo siento, estoy teniendo problemas de conexión. Intenta más tarde.", nil
	}

	msg := resp.Choices[0].Message

	if len(msg.ToolCalls) > 0 {
		openaiMsgs = append(openaiMsgs, msg)

		for _, toolCall := range msg.ToolCalls {
			if toolCall.Type == openai.ToolTypeFunction {
				toolOutput, err := s.callFunction(ctx, userID, toolCall.Function.Name, toolCall.Function.Arguments)
				if err != nil {
					toolOutput = fmt.Sprintf("Error executing tool: %v", err)
				}

				openaiMsgs = append(openaiMsgs, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    toolOutput,
					ToolCallID: toolCall.ID,
				})
			}
		}

		resp, err = s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: openaiMsgs,
			Tools:    tools,
		})
		if err != nil {
			s.logger.Errorw("Error calling OpenAI after tools", "error", err)
			return "Lo siento, hubo un error procesando la solicitud.", nil
		}
		msg = resp.Choices[0].Message
	}

	botContent := msg.Content

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

func (s *LLMService) GetChatHistory(ctx context.Context, userID uuid.UUID, fq pagination.PaginatedFeedQuery) ([]llmdomain.Message, int64, error) {
	convo, err := s.store.LLM.GetActiveConversation(ctx, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("error searching for active conversation: %v", err)
	}

	if convo == nil {
		return []llmdomain.Message{}, 0, nil
	}

	history, total, err := s.store.LLM.GetConversationHistory(ctx, convo.ConversationID, fq)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting history: %v", err)
	}

	return history, total, nil
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
