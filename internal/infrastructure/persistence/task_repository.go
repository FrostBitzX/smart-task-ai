package persistence

import (
	"context"
	"time"

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

func (r *taskRepository) CountTasksByStatus(ctx context.Context, nodeID uuid.UUID) ([]tasks.StatusCount, error) {
	var results []tasks.StatusCount
	err := r.db.WithContext(ctx).
		Model(&entity.Task{}).
		Select("status, COUNT(*) as count").
		Where("node_id = ?", nodeID).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *taskRepository) ListUnscheduledTasks(ctx context.Context, nodeID uuid.UUID) ([]*entity.Task, error) {
	var tasks []*entity.Task
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("node_id = ?", nodeID).
		Where("status IN (?)", []string{"todo", "in_progress"}).
		Where("start_datetime IS NULL").
		Where("end_datetime IS NULL").
		Order("updated_at DESC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) ListTodayTasks(ctx context.Context, nodeID uuid.UUID, today time.Time) ([]*entity.Task, error) {
	// Get start of day (00:00:00)
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	// Get end of day (23:59:59)
	endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())

	var tasks []*entity.Task
	err := r.db.WithContext(ctx).
		Preload("Project").
		Where("node_id = ?", nodeID).
		Where("start_datetime >= ? AND start_datetime <= ?", startOfDay, endOfDay).
		Order("start_datetime ASC").
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
