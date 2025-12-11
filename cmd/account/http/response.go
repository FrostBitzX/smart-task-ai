package http

import "github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"

// AccountResponse is the API response shape for an account payload.
type AccountResponse struct {
	Username string `json:"username"`
}

// toAPIAccountResponse converts an accounts.Account DTO into an AccountResponse
// struct that will be embedded in the success response Data field.
func toAPIAccountResponse(acc accounts.Account) AccountResponse {
	return AccountResponse{
		Username: acc.Username,
	}
}
