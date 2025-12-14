package accounts

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
)

type AccountRepository interface {
	Create(ctx context.Context, acc *entity.Account) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
