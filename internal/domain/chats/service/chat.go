package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/chats"
	projectEntity "github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
	projectSvc "github.com/FrostBitzX/smart-task-ai/internal/domain/projects/service"
	taskSvc "github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/service"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/groq"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"
)

// Error codes for chat service
const (
	ErrCodeGroqUnavailable = "GROQ_UNAVAILABLE"
	ErrCodeGroqTimeout     = "GROQ_TIMEOUT"
	ErrCodeRateLimited     = "RATE_LIMITED"
	ErrCodeGroqAuthError   = "GROQ_AUTH_ERROR"
	ErrCodeProjectNotFound = "PROJECT_NOT_FOUND"
	ErrCodeInvalidMessage  = "INVALID_MESSAGE"
)

// Pre-compiled regex for extracting JSON from markdown code blocks
var codeBlockRegex = regexp.MustCompile("```(?:json)?\\s*([\\s\\S]*?)```")

// ChatService defines the interface for chat operations
type ChatService interface {
	SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
	SendMessageStream(ctx context.Context, req *SendMessageRequest) (<-chan groq.StreamChunk, error)
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	ProjectID      uuid.UUID
	AccountID      uuid.UUID
	Content        string
	SessionHistory []groq.ChatMessage
}

// SendMessageResponse represents the response from sending a message
type SendMessageResponse struct {
	Message     string
	TaskActions []TaskAction
}

// TaskAction represents an action performed on a task
type TaskAction struct {
	Type   string
	TaskID uuid.UUID
	Name   string
}

// chatService implements the ChatService interface
type chatService struct {
	groqClient       groq.GroqClient
	taskService      *taskSvc.TaskService
	projectService   *projectSvc.ProjectService
	promptBuilder    PromptBuilder
	functionExecutor *FunctionExecutor
}

// NewChatService creates a new chat service
func NewChatService(
	groqClient groq.GroqClient,
	taskService *taskSvc.TaskService,
	projectService *projectSvc.ProjectService,
) ChatService {
	return &chatService{
		groqClient:       groqClient,
		taskService:      taskService,
		projectService:   projectService,
		promptBuilder:    NewPromptBuilder(),
		functionExecutor: NewFunctionExecutor(taskService),
	}
}

// SendMessage sends a message to the AI and returns the response
func (s *chatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	project, err := s.projectService.GetProjectByID(ctx, req.ProjectID)
	if err != nil {
		return nil, s.handleProjectError(err)
	}

	tasks, err := s.taskService.ListTasksByProject(ctx, req.ProjectID)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to get tasks", "GET_TASKS_ERROR", err)
	}

	aiConfig := s.getAIConfig(project)
	systemPrompt := s.promptBuilder.BuildSystemPrompt(aiConfig, tasks)
	messages := s.buildMessages(systemPrompt, req.SessionHistory, req.Content)
	groqReq := groq.NewDefaultRequest(messages)

	resp, err := s.groqClient.SendChatCompletion(ctx, groqReq)
	if err != nil {
		return nil, s.handleGroqError(err)
	}

	if len(resp.Choices) == 0 {
		return nil, apperror.NewInternalServerError("no response from AI", "EMPTY_RESPONSE", nil)
	}

	aiResponse := resp.Choices[0].Message.Content

	// Try to parse AI response as a function call
	actionRequest := s.parseActionFromResponse(aiResponse)
	if actionRequest != nil {
		// Execute the function call
		result := s.functionExecutor.ExecuteFunction(ctx, req.ProjectID, &FunctionCall{
			Name:      actionRequest.Action,
			Arguments: actionRequest.Params,
		})

		var taskActions []TaskAction
		if result.Action != nil {
			taskActions = append(taskActions, *result.Action)
		}

		// Generate a human-readable response based on the result
		responseMessage := s.generateActionResponse(actionRequest.Action, result)

		return &SendMessageResponse{
			Message:     responseMessage,
			TaskActions: taskActions,
		}, nil
	}

	return &SendMessageResponse{
		Message:     aiResponse,
		TaskActions: []TaskAction{},
	}, nil
}

// ActionRequest represents a parsed action from AI response
type ActionRequest struct {
	Action string          `json:"action"`
	Params json.RawMessage `json:"params"`
}

// parseActionFromResponse tries to parse an action request from AI response
func (s *chatService) parseActionFromResponse(response string) *ActionRequest {
	// Try to extract JSON from the response
	jsonStr := extractJSON(response)
	if jsonStr == "" {
		return nil
	}

	var actionReq ActionRequest
	if err := json.Unmarshal([]byte(jsonStr), &actionReq); err != nil {
		return nil
	}

	// Validate action name
	validActions := map[string]bool{
		"create_task": true,
		"update_task": true,
		"delete_task": true,
		"get_task":    true,
		"list_tasks":  true,
	}

	if !validActions[actionReq.Action] {
		return nil
	}

	return &actionReq
}

// extractJSON extracts JSON object from a string (handles markdown code blocks)
func extractJSON(s string) string {
	s = strings.TrimSpace(s)

	// Try to extract from markdown code block (using pre-compiled regex)
	matches := codeBlockRegex.FindStringSubmatch(s)
	if len(matches) > 1 {
		s = strings.TrimSpace(matches[1])
	}

	// Find JSON object boundaries
	start := strings.Index(s, "{")
	if start == -1 {
		return ""
	}

	// Find matching closing brace
	depth := 0
	for i := start; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}

	return ""
}

// generateActionResponse generates a human-readable response for an action result
func (s *chatService) generateActionResponse(action string, result *FunctionResult) string {
	if !result.Success {
		return "ขออภัย ไม่สามารถดำเนินการได้: " + result.Error
	}

	switch action {
	case "create_task":
		if data, ok := result.Data.(map[string]interface{}); ok {
			name := data["name"]
			return "สร้าง task \"" + name.(string) + "\" เรียบร้อยแล้ว ✅"
		}
	case "update_task":
		if data, ok := result.Data.(map[string]interface{}); ok {
			name := data["name"]
			return "อัพเดท task \"" + name.(string) + "\" เรียบร้อยแล้ว ✅"
		}
	case "delete_task":
		return "ลบ task เรียบร้อยแล้ว ✅"
	case "get_task":
		if data, ok := result.Data.(map[string]interface{}); ok {
			return formatTaskDetails(data)
		}
	case "list_tasks":
		if data, ok := result.Data.(map[string]interface{}); ok {
			return formatTaskList(data)
		}
	}

	return "ดำเนินการเรียบร้อยแล้ว ✅"
}

func formatTaskDetails(data map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString("📋 รายละเอียด Task:\n")
	sb.WriteString("- ชื่อ: " + data["name"].(string) + "\n")
	sb.WriteString("- สถานะ: " + data["status"].(string) + "\n")
	sb.WriteString("- Priority: " + data["priority"].(string) + "\n")
	if desc, ok := data["description"]; ok {
		sb.WriteString("- รายละเอียด: " + desc.(string) + "\n")
	}
	if start, ok := data["start_datetime"]; ok {
		sb.WriteString("- เริ่ม: " + start.(string) + "\n")
	}
	if end, ok := data["end_datetime"]; ok {
		sb.WriteString("- สิ้นสุด: " + end.(string) + "\n")
	}
	return sb.String()
}

func formatTaskList(data map[string]interface{}) string {
	tasks, ok := data["tasks"].([]map[string]interface{})
	if !ok {
		return "ไม่พบ tasks"
	}

	if len(tasks) == 0 {
		return "📋 ไม่มี task ใน project นี้"
	}

	var sb strings.Builder
	sb.WriteString("📋 รายการ Tasks:\n")
	for i, t := range tasks {
		name := t["name"].(string)
		status := t["status"].(string)
		priority := t["priority"].(string)
		sb.WriteString(strings.Repeat("-", 30) + "\n")
		sb.WriteString(strings.Repeat(" ", 2) + string(rune('1'+i)) + ". " + name + "\n")
		sb.WriteString("     สถานะ: " + status + " | Priority: " + priority + "\n")
	}
	return sb.String()
}

// SendMessageStream sends a message to the AI and returns a streaming response
func (s *chatService) SendMessageStream(ctx context.Context, req *SendMessageRequest) (<-chan groq.StreamChunk, error) {
	if err := s.validateRequest(req); err != nil {
		return nil, err
	}

	project, err := s.projectService.GetProjectByID(ctx, req.ProjectID)
	if err != nil {
		return nil, s.handleProjectError(err)
	}

	tasks, err := s.taskService.ListTasksByProject(ctx, req.ProjectID)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to get tasks", "GET_TASKS_ERROR", err)
	}

	aiConfig := s.getAIConfig(project)
	systemPrompt := s.promptBuilder.BuildSystemPrompt(aiConfig, tasks)
	messages := s.buildMessages(systemPrompt, req.SessionHistory, req.Content)

	groqReq := groq.NewDefaultRequest(messages)
	groqReq.Stream = true

	return s.groqClient.SendChatCompletionStream(ctx, groqReq)
}

func (s *chatService) validateRequest(req *SendMessageRequest) error {
	if req == nil {
		return apperror.NewBadRequestError("request is required", ErrCodeInvalidMessage, nil)
	}
	if strings.TrimSpace(req.Content) == "" {
		return apperror.NewBadRequestError("message content is required", ErrCodeInvalidMessage, nil)
	}
	if req.ProjectID == uuid.Nil {
		return apperror.NewBadRequestError("project ID is required", ErrCodeInvalidMessage, nil)
	}
	return nil
}

func (s *chatService) handleProjectError(err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		if appErr.Code == "PROJECT_NOT_FOUND" {
			return apperror.NewNotFoundError("project not found", ErrCodeProjectNotFound, err)
		}
	}
	return apperror.NewInternalServerError("failed to get project", "GET_PROJECT_ERROR", err)
}

func (s *chatService) handleGroqError(err error) error {
	errMsg := err.Error()

	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
		return apperror.NewTimeoutError("AI response timeout", ErrCodeGroqTimeout, err)
	}

	if strings.Contains(errMsg, "status 429") || strings.Contains(errMsg, "Too many requests") {
		return apperror.NewTooManyRequestsError("Too many requests", ErrCodeRateLimited, err)
	}

	if strings.Contains(errMsg, "status 401") || strings.Contains(errMsg, "Unauthorized") {
		return apperror.NewInternalServerError("AI service configuration error", ErrCodeGroqAuthError, err)
	}

	if strings.Contains(errMsg, "status 503") || strings.Contains(errMsg, "unavailable") {
		return &apperror.AppError{
			Status:   http.StatusServiceUnavailable,
			Code:     ErrCodeGroqUnavailable,
			Message:  "AI service temporarily unavailable",
			RawError: err,
		}
	}

	return &apperror.AppError{
		Status:   http.StatusServiceUnavailable,
		Code:     ErrCodeGroqUnavailable,
		Message:  "AI service temporarily unavailable",
		RawError: err,
	}
}

func (s *chatService) getAIConfig(project *projectEntity.Project) *chats.AIConfig {
	if len(project.Config) == 0 {
		return &chats.DefaultAIConfig
	}

	var configWrapper struct {
		AIConfig *chats.AIConfig `json:"ai_config"`
	}

	if err := json.Unmarshal(project.Config, &configWrapper); err != nil {
		return &chats.DefaultAIConfig
	}

	if configWrapper.AIConfig == nil {
		return &chats.DefaultAIConfig
	}

	config := configWrapper.AIConfig
	if config.ChatStyle == "" {
		config.ChatStyle = chats.DefaultAIConfig.ChatStyle
	}
	if len(config.DomainKnowledge) == 0 {
		config.DomainKnowledge = chats.DefaultAIConfig.DomainKnowledge
	}
	if config.Language == "" {
		config.Language = chats.DefaultAIConfig.Language
	}

	return config
}

func (s *chatService) buildMessages(systemPrompt string, sessionHistory []groq.ChatMessage, userContent string) []groq.ChatMessage {
	messages := make([]groq.ChatMessage, 0, len(sessionHistory)+2)

	messages = append(messages, groq.ChatMessage{
		Role:    "system",
		Content: systemPrompt,
	})

	messages = append(messages, sessionHistory...)

	messages = append(messages, groq.ChatMessage{
		Role:    "user",
		Content: userContent,
	})

	return messages
}
