package persistence

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/google/uuid"
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

func (r *taskRepository) GetTaskByID(ctx context.Context, taskID uuid.UUID, nodeID uuid.UUID) (*entity.Task, error) {
	var task entity.Task
	err := r.db.WithContext(ctx).
		Where("id = ? AND node_id = ?", taskID, nodeID).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) ListTasksByProject(ctx context.Context, projectID uuid.UUID, nodeID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND node_id = ?", projectID, nodeID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) CountTasksByProject(ctx context.Context, projectID uuid.UUID, nodeID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entity.Task{}).
		Where("project_id = ? AND node_id = ?", projectID, nodeID).
		Count(&count).Error
	return count, err
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task, nodeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.Task{}).
		Where("id = ? AND node_id = ?", task.ID, nodeID).
		Updates(task).Error
}

func (r *taskRepository) DeleteTask(ctx context.Context, taskID uuid.UUID, nodeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND node_id = ?", taskID, nodeID).
		Delete(&entity.Task{}).Error
}
