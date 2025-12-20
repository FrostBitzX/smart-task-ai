package profiles

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, prof *entity.Profile) error
}
