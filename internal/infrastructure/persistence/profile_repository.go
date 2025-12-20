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
