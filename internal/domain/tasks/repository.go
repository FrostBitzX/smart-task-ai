package tasks

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *entity.Task) error
}
