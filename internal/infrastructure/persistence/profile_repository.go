package persistence

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
	"gorm.io/gorm"
)

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) profiles.ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) CreateProfile(ctx context.Context, prof *entity.Profile) error {
	return r.db.WithContext(ctx).Create(prof).Error
}

func (r *profileRepository) CheckAndGetProfile(ctx context.Context, accountID string) (*entity.Profile, error) {
	var profile entity.Profile
	err := r.db.WithContext(ctx).
		Select("id, account_id, first_name, last_name, nickname, avatar_path, state, created_at, updated_at").
		Where("account_id = ?", accountID).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) UpdateProfile(ctx context.Context, prof *entity.Profile) error {
	return r.db.WithContext(ctx).Save(prof).Error
}
