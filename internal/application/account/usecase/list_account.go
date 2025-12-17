package usecase

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/service"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
)

type ListAccountUseCase struct {
	accountService *service.AccountService
	logger         logger.Logger
}

func NewListAccountUseCase(svc *service.AccountService, l logger.Logger) *ListAccountUseCase {
	return &ListAccountUseCase{
		accountService: svc,
		logger:         l,
	}
}

func (uc *ListAccountUseCase) Execute(ctx context.Context, req *account.ListAccountsRequest) (*account.ListAccountsResponse, error) {
	if req == nil {
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// Set pagination defaults
	var limit int
	var offset int
	if req.Limit == nil {
		limit = 10
	} else {
		limit = *req.Limit
	}

	if req.Offset == nil {
		offset = 0
	} else {
		offset = *req.Offset
	}

	// Get accounts from service
	accounts, total, err := uc.accountService.ListAccounts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert entities to DTOs
	accountDTOs := make([]account.AccountDTO, len(accounts))
	for i, acc := range accounts {
		accountDTOs[i] = account.AccountDTO{
			ID:       utils.ShortUUIDWithPrefix(acc.ID, entity.AccountIDPrefix),
			Username: acc.Username,
			Email:    acc.Email,
			Status:   acc.State,
		}
	}

	// Calculate pagination info
	hasMore := offset+limit < total

	// Build response
	response := &account.ListAccountsResponse{
		Items: accountDTOs,
		Pagination: account.PaginationDTO{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}

	return response, nil
}
