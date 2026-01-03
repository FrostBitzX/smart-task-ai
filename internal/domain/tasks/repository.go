package tasks

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *entity.Task) error
	GetTaskByID(ctx context.Context, taskID uuid.UUID) (*entity.Task, error)
	ListTasksByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Task, error)
}
