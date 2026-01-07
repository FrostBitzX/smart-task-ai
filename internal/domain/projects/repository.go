//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=../../mocks/project_repository.go -package=mocks
package projects

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
	"github.com/google/uuid"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, proj *entity.Project) error
	GetProjectByID(ctx context.Context, projectID uuid.UUID) (*entity.Project, error)
	ListProjectByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Project, int, error)
}
