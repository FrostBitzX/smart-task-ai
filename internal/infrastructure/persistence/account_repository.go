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

func (r *accountRepository) ExistsAccount(ctx context.Context, username, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Account{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error

	return count > 0, err
}

func (r *accountRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	var account entity.Account

	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&account).Error

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *accountRepository) ListAccounts(ctx context.Context, limit, offset int) ([]*entity.Account, int, error) {
	var accounts []*entity.Account
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&entity.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&accounts).Error

	return accounts, int(total), err
}

func (r *accountRepository) GetAccount(ctx context.Context, id string) (*entity.Account, error) {
	var account entity.Account
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&account).Error

	if err != nil {
		return nil, err
	}

	return &account, nil
}
