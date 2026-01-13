package chat

import "encoding/json"

// SendMessageRequestDTO represents the request to send a message to the AI assistant
type SendMessageRequestDTO struct {
	ProjectID      string       `json:"project_id,omitempty"` // Set from URL parameter
	Content        string       `json:"content" validate:"required"`
	SessionHistory []MessageDTO `json:"session_history,omitempty"`
}

// MessageDTO represents a message in the conversation history
type MessageDTO struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SendMessageResponseDTO represents the response from the AI assistant
type SendMessageResponseDTO struct {
	Message     json.RawMessage `json:"message"`
	TaskActions []TaskActionDTO `json:"task_actions,omitempty"`
}

// TaskActionDTO represents an action performed on a task
type TaskActionDTO struct {
	Type   string `json:"type"`
	TaskID string `json:"task_id"`
	Name   string `json:"name"`
}
