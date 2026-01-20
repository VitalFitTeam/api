package llmhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary		Chat with AI Assistant
// @Description	Sends a message to the AI assistant and receives a response.
// @Tags			LLM
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			payload	body		ChatRequest				true	"Chat message payload"
// @Success		200		{object}	ChatResponse			"AI response"
// @Failure		400		{object}	object{error=string}	"Bad Request"
// @Failure		500		{object}	object{error=string}	"Internal Server Error"
// @Router			/llm/chat [post]
func (h *LLMHandlers) ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.services.LogErrors.BadRequestResponse(c, err)
		return
	}

	user := h.services.UserServices.GetUserFromContext(c)

	botResponse, err := h.services.LLM.ProcessUserMessage(c.Request.Context(), user.UserID, req.Message)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, ChatResponse{Response: botResponse})
}

// @Summary		Get chat history
// @Description	Retrieves the conversation history for the current user.
// @Tags			LLM
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{array}		llmdomain.Message		"Chat history"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/llm/history [get]
func (h *LLMHandlers) GetHistoryHandler(c *gin.Context) {
	// 1. Obtener usuario del contexto (Tu helper)
	user := h.services.UserServices.GetUserFromContext(c)

	// 2. Llamar al servicio
	history, err := h.services.LLM.GetChatHistory(c.Request.Context(), user.UserID)
	if err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	// 3. Responder
	c.JSON(http.StatusOK, history)
}

// @Summary		Reset conversation
// @Description	Resets the active conversation for the current user.
// @Tags			LLM
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	map[string]string		"Success message"
// @Failure		500	{object}	object{error=string}	"Internal Server Error"
// @Router			/llm/reset [post]
func (h *LLMHandlers) ResetChatHandler(c *gin.Context) {
	// 1. Obtener usuario
	user := h.services.UserServices.GetUserFromContext(c)

	// 2. Llamar al servicio
	if err := h.services.LLM.ResetConversation(c.Request.Context(), user.UserID); err != nil {
		h.services.LogErrors.InternalServerError(c, err)
		return
	}

	// 3. Responder
	c.JSON(http.StatusOK, gin.H{"message": "Conversation reset successfully"})
}
