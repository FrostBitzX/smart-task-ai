package usecase

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/service"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
)

type LoginUseCase struct {
	accountService *service.AccountService
	logger         logger.Logger
}

func NewLoginUseCase(accSvc *service.AccountService, l logger.Logger) *LoginUseCase {
	return &LoginUseCase{
		accountService: accSvc,
		logger:         l,
	}
}

func (u *LoginUseCase) Execute(req *account.LoginRequest) (*account.LoginResponse, error) {
	token, err := u.accountService.Login(context.Background(), req)

	if err != nil {
		return nil, err
	}

	res := &account.LoginResponse{
		Token: token,
	}

	return res, nil
}
