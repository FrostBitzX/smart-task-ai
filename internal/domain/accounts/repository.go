package accounts

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *entity.Account) error
	ExistsAccount(ctx context.Context, username, email string) (bool, error)
}
