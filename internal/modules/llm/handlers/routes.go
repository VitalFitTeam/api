package llmhandlers

import appservices "github.com/vitalfit/api/internal/app/services"

type LLMHandlersInterface interface{}

type LLMHandlers struct {
	appservices appservices.Services
}

func NewLLMHandlers(appservices appservices.Services) *LLMHandlers {
	return &LLMHandlers{
		appservices: appservices,
	}
}
