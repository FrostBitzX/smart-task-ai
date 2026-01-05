package usecase

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/application/chat"
	chatSvc "github.com/FrostBitzX/smart-task-ai/internal/domain/chats/service"
	projectEntity "github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/groq"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"
)

// SendMessageUseCase handles sending messages to the AI assistant
type SendMessageUseCase struct {
	chatService chatSvc.ChatService
	logger      logger.Logger
}

// NewSendMessageUseCase creates a new SendMessageUseCase
func NewSendMessageUseCase(cs chatSvc.ChatService, l logger.Logger) *SendMessageUseCase {
	return &SendMessageUseCase{
		chatService: cs,
		logger:      l,
	}
}

// ExecuteStream sends a message to the AI assistant and returns a streaming response
func (uc *SendMessageUseCase) ExecuteStream(ctx context.Context, accountID string, req *chat.SendMessageRequestDTO) (<-chan groq.StreamChunk, error) {
	if req == nil {
		return nil, apperror.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	parsedProjectID, err := utils.ParseID(req.ProjectID, projectEntity.ProjectIDPrefix)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid project ID format", "INVALID_PROJECT_ID", err)
	}

	parsedAccountID, err := uuid.Parse(accountID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid account ID format", "INVALID_ACCOUNT_ID", err)
	}

	return uc.chatService.SendMessageStream(ctx, &chatSvc.SendMessageRequest{
		ProjectID:      parsedProjectID,
		AccountID:      parsedAccountID,
		Content:        req.Content,
		SessionHistory: convertSessionHistory(req.SessionHistory),
	})
}

// convertSessionHistory converts DTOs to domain messages
func convertSessionHistory(history []chat.MessageDTO) []groq.ChatMessage {
	if len(history) == 0 {
		return nil
	}

	messages := make([]groq.ChatMessage, len(history))
	for i, msg := range history {
		messages[i] = groq.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return messages
}
