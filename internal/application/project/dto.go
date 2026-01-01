package project

import "encoding/json"

type CreateProjectRequest struct {
	AccountID string          `json:"account_id"`
	Name      string          `json:"name" validate:"min=0,max=50"`
	Config    json.RawMessage `json:"config"`
}

type CreateProjectResponse struct {
	ProjectID string `json:"project_id"`
}
