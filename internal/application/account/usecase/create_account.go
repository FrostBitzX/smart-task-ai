package usecase

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/service"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
)

type AccountUseCase struct {
	accountService *service.AccountService
	logger         logger.Logger
}

func NewAccountUseCase(svc *service.AccountService, l logger.Logger) *AccountUseCase {
	return &AccountUseCase{
		accountService: svc,
		logger:         l,
	}
}

func (uc *AccountUseCase) Execute(req *account.CreateAccountRequest) (*account.CreateAccountResponse, error) {
	acc, err := uc.accountService.CreateAccount(context.Background(), req)

	if err != nil {
		return nil, err
	}

	// Convert UUID to string with prefix
	accountID := utils.ShortUUIDWithPrefix(acc.ID, entity.AccountIDPrefix)

	res := &account.CreateAccountResponse{
		AccountID: accountID,
	}
	return res, nil
}
