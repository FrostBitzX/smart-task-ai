package service

import (
	"context"
	"errors"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/google/uuid"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
)

type AccountService struct {
	repo accounts.AccountRepository
}

func NewAccountService(repo accounts.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

// [TODO]: not return entire entity, use model
func (s *AccountService) CreateAccount(ctx context.Context, req *account.CreateAccountRequest) (*entity.Account, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("invalid input")
	}

	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// create domain entity
	acc := &entity.Account{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // [TODO]: hash password
		// Password: hash(req.Password) // แนะนำทำ hash ที่นี่หรือ infra
	}

	// persist
	if err := s.repo.Create(ctx, acc); err != nil {
		return nil, err
	}

	return acc, nil
}
