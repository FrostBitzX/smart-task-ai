package service

import (
	"context"
	"os"
	"testing"

	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAccountService_CreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	svc := NewAccountService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		req := &account.CreateAccountRequest{
			Username:        "testuser",
			Email:           "test@example.com",
			Password:        "password123",
			ConfirmPassword: "password123",
		}

		mockRepo.EXPECT().
			ExistsAccount(ctx, req.Username, req.Email).
			Return(false, nil).
			Times(1)
		mockRepo.EXPECT().
			CreateAccount(ctx, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := svc.CreateAccount(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Username, res.Username)
	})

	t.Run("account already exists", func(t *testing.T) {
		req := &account.CreateAccountRequest{
			Username: "existinguser",
			Email:    "existing@example.com",
		}

		mockRepo.EXPECT().
			ExistsAccount(ctx, req.Username, req.Email).
			Return(true, nil).
			Times(1)

		res, err := svc.CreateAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("password mismatch", func(t *testing.T) {
		req := &account.CreateAccountRequest{
			Username:        "user",
			Email:           "email@example.com",
			Password:        "pass1",
			ConfirmPassword: "pass2",
		}

		mockRepo.EXPECT().
			ExistsAccount(ctx, req.Username, req.Email).
			Return(false, nil).
			Times(1)

		res, err := svc.CreateAccount(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "password and confirm password does not match")
	})
}

func TestAccountService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	svc := NewAccountService(mockRepo)
	ctx := context.Background()

	os.Setenv("JWT_SECRET", "testsecret")
	defer os.Unsetenv("JWT_SECRET")

	t.Run("success", func(t *testing.T) {
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		acc := &entity.Account{
			ID:       uuid.New(),
			Username: "testuser",
			Email:    "test@example.com",
			Password: string(hashedPassword),
		}

		req := &account.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		mockRepo.EXPECT().
			GetByUsername(ctx, req.Username).
			Return(acc, nil).
			Times(1)

		token, err := svc.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		acc := &entity.Account{
			Username: "testuser",
			Password: "wrongpassword", // This is not a BCrypt hash so comparison will fail
		}

		req := &account.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		mockRepo.EXPECT().
			GetByUsername(ctx, req.Username).
			Return(acc, nil).
			Times(1)

		token, err := svc.Login(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid username or password")
	})
}
