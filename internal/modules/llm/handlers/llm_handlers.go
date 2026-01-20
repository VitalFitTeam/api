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
