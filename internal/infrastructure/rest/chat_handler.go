package rest

import (
	"github.com/FrostBitzX/smart-task-ai/internal/application/chat"
	"github.com/FrostBitzX/smart-task-ai/internal/application/chat/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
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

	resp, err := h.sendMessageUC.Execute(c.Context(), accountID, req)
	if err != nil {
		h.logger.Error("Failed to send message", map[string]interface{}{
			"error":      err.Error(),
			"account_id": accountID,
			"project_id": req.ProjectID,
		})
		return responses.Error(c, err)
	}

	return responses.Success(c, resp, "Message sent successfully")
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
