package projects

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, proj *entity.Project) error
}
