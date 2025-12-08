package accounts

import "github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"

type AccountRepository interface {
	FindByEmail(email string) (*entity.Account, error)
}
