package service

import "github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"

type AccountService struct {
	repo accounts.AccountRepository
}

func NewAccountService(repo accounts.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}
