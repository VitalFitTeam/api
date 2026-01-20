package llmservices

import (
	"github.com/sashabaranov/go-openai"
	"github.com/vitalfit/api/internal/store"
	"go.uber.org/zap"
)

type LLMService struct {
	client *openai.Client
	store  store.Storage
	logger *zap.SugaredLogger
}

func NewLLMService(apiKey string, store store.Storage, logger *zap.SugaredLogger) *LLMService {
	client := openai.NewClient(apiKey)
	return &LLMService{
		client: client,
		store:  store,
		logger: logger,
	}
}
