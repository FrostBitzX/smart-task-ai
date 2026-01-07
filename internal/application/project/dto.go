package project

import (
	"encoding/json"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/common"
)

type CreateProjectRequest struct {
	AccountID string          `json:"account_id"`
	Name      string          `json:"name" validate:"min=0,max=50"`
	Config    json.RawMessage `json:"config"`
}

type ListProjectRequest struct {
	AccountID string `query:"account_id" validate:"omitempty"`
	Limit     *int   `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset    *int   `query:"offset" validate:"omitempty,min=0"`
}

type CreateProjectResponse struct {
	ProjectID string `json:"project_id"`
}

type ProjectResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Config    json.RawMessage `json:"config"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ListProjectResponse struct {
	Items      []ProjectResponse `json:"items"`
	Pagination common.Pagination `json:"pagination"`
}
