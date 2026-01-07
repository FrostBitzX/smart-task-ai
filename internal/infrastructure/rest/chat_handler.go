package rest

import (
	"bufio"
	"fmt"

	"github.com/FrostBitzX/smart-task-ai/internal/application/chat"
	"github.com/FrostBitzX/smart-task-ai/internal/application/chat/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// ChatHandler handles chat-related HTTP requests
type ChatHandler struct {
	sendMessageUC *usecase.SendMessageUseCase
	logger        logger.Logger
}

// NewChatHandler creates a new ChatHandler
func NewChatHandler(sendMessageUC *usecase.SendMessageUseCase, l logger.Logger) *ChatHandler {
	return &ChatHandler{
		sendMessageUC: sendMessageUC,
		logger:        l,
	}
}

// SendMessage handles POST /api/:projectId/chat endpoint
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	// Get projectId from URL parameter
	projectID := c.Params("projectId")
	if projectID == "" {
		return responses.Error(c, apperror.NewBadRequestError("project ID is required", "INVALID_PROJECT_ID", nil))
	}

	// Check if streaming is requested via query param (default: false)
	stream := c.Query("stream", "false") == "true"

	req, err := requests.ParseAndValidate[chat.SendMessageRequestDTO](c)
	if err != nil {
		h.logger.Warn("Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return responses.Error(c, err)
	}

	// Set projectID from URL parameter
	req.ProjectID = projectID

	// Get account ID from JWT claims
	accountID, err := h.getAccountIDFromContext(c)
	if err != nil {
		return responses.Error(c, err)
	}

	// Use streaming only when explicitly requested with ?stream=true
	if stream {
		return h.handleStreamResponse(c, accountID, req)
	}
	return h.handleNonStreamResponse(c, accountID, req)
}

// handleNonStreamResponse handles non-streaming JSON response
func (h *ChatHandler) handleNonStreamResponse(c *fiber.Ctx, accountID string, req *chat.SendMessageRequestDTO) error {
	resp, err := h.sendMessageUC.Execute(c.Context(), accountID, req)
	if err != nil {
		h.logger.Error("Failed to send message", map[string]interface{}{
			"error":      err.Error(),
			"account_id": accountID,
			"project_id": req.ProjectID,
		})
		return responses.Error(c, err)
	}

	return responses.Success(c, resp)
}

// handleStreamResponse handles streaming SSE response
func (h *ChatHandler) handleStreamResponse(c *fiber.Ctx, accountID string, req *chat.SendMessageRequestDTO) error {
	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Get streaming channel
	streamChan, err := h.sendMessageUC.ExecuteStream(c.Context(), accountID, req)
	if err != nil {
		h.logger.Error("Failed to start stream", map[string]interface{}{
			"error":      err.Error(),
			"account_id": accountID,
			"project_id": req.ProjectID,
		})
		return responses.Error(c, err)
	}

	// Stream response using Fiber's streaming
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for chunk := range streamChan {
			if chunk.Error != nil {
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", chunk.Error.Error())
				w.Flush()
				return
			}

			if chunk.Done {
				fmt.Fprintf(w, "event: done\ndata: [DONE]\n\n")
				w.Flush()
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", chunk.Content)
			w.Flush()
		}
	}))

	return nil
}

// getAccountIDFromContext extracts account ID from JWT claims in context
func (h *ChatHandler) getAccountIDFromContext(c *fiber.Ctx) (string, error) {
	claims, ok := c.Locals("jwt_claims").(map[string]interface{})
	if !ok {
		return "", apperror.NewUnauthorizedError("authentication required", "UNAUTHORIZED", nil)
	}

	accountID, ok := claims["AccountId"].(string)
	if !ok || accountID == "" {
		return "", apperror.NewUnauthorizedError("invalid token claims", "INVALID_TOKEN_CLAIMS", nil)
	}

	return accountID, nil
}
