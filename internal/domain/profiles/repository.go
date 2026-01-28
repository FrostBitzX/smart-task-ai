//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=../../mocks/profile_repository.go -package=mocks
package profiles

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
)

type ProfileRepository interface {
	CreateProfile(ctx context.Context, prof *entity.Profile) error
	GetProfile(ctx context.Context, accountID string, nodeID string) (*entity.Profile, error)
	UpdateProfile(ctx context.Context, prof *entity.Profile, nodeID string) error
}
