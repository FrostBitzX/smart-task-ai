package service

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
)

type AccountService struct {
	repo accounts.AccountRepository
}

func NewAccountService(repo accounts.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *account.CreateAccountRequest) (*entity.Account, error) {
	// Check if username or email already exists
	exists, err := s.repo.ExistsAccount(ctx, req.Username, req.Email)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to check account existence", "CHECK_ACCOUNT_EXISTS_ERROR", err)
	}
	if exists {
		return nil, apperrors.NewBadRequestError("username or email already exists", "USERNAME_OR_EMAIL_EXISTS", nil)
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to hash password", "HASH_PASSWORD_ERROR", err)
	}

	// create domain entity
	acc := &entity.Account{
		ID:       uuid.New(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// persist account to database
	if err := s.repo.CreateAccount(ctx, acc); err != nil {
		return nil, apperrors.NewInternalServerError("failed to create account", "CREATE_ACCOUNT_ERROR", err)
	}

	return acc, nil
}
