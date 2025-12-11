package http

import (
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/service"
)

// AccountHandlerInterface defines the interface for account HTTP handlers
type AccountHandlerInterface interface {
	StrictServerInterface
}

// AccountHandler handles HTTP requests for account operations
type AccountHandler struct {
	AccountService *service.AccountService
}

// NewAccountHandler creates a new instance of AccountHandler
func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{
		AccountService: accountService,
	}
}

// Ensure AccountHandler implements StrictServerInterface
var _ StrictServerInterface = (*AccountHandler)(nil)
