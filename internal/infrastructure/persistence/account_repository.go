package persistence

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/entity"
	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) accounts.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, acc *entity.Account) error {
	return r.db.WithContext(ctx).Create(acc).Error
}

func (r *accountRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Account{}).
		Where("username = ?", username).
		Count(&count).Error

	return count > 0, err
}
