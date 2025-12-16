package service

import (
	"context"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/golang-jwt/jwt/v5"
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
	now := time.Now()
	acc := &entity.Account{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		State:     "active",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// persist account to database
	if err := s.repo.CreateAccount(ctx, acc); err != nil {
		return nil, apperrors.NewInternalServerError("failed to create account", "CREATE_ACCOUNT_ERROR", err)
	}

	return acc, nil
}

func (s *AccountService) Login(ctx context.Context, req *account.LoginRequest) (string, error) {
	acc, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", apperrors.NewBadRequestError("user does not exist", "LOGIN_ERROR", nil)
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(req.Password)); err != nil {
		return "", apperrors.NewBadRequestError("invalid username or password", "LOGIN_ERROR", nil)
	}

	claims := jwt.MapClaims{
		"name": acc.Username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", apperrors.NewInternalServerError(
			"failed to sign jwt",
			"JWT_SIGN_ERROR",
			err,
		)
	}

	return t, nil
}
