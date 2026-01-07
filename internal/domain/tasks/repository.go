//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=../../mocks/task_repository.go -package=mocks
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
	CountTasksByProject(ctx context.Context, projectID uuid.UUID) (int64, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
}
