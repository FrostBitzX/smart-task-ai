//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=../../mocks/project_repository.go -package=mocks
package projects

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, proj *entity.Project) error
}
