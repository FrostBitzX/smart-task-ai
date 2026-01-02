package persistence

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) tasks.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}
