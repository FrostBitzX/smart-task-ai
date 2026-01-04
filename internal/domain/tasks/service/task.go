package service

import (
	"context"
	"errors"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
)

type TaskService struct {
	repo tasks.TaskRepository
}

func NewTaskService(repo tasks.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, projectID uuid.UUID, req *task.CreateTaskRequest) (*entity.Task, error) {
	if req == nil {
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	if req.StartDateTime != nil && req.EndDateTime != nil {
		if *req.StartDateTime == *req.EndDateTime {
			return nil, apperrors.NewBadRequestError("start_datetime and end_datetime cannot be the same", "INVALID_REQUEST", nil)
		}
		if *req.EndDateTime < *req.StartDateTime {
			return nil, apperrors.NewBadRequestError("end_datetime must be greater than start_datetime", "INVALID_REQUEST", nil)
		}
	}

	// create domain entity
	now := time.Now()
	task := &entity.Task{
		ID:            uuid.New(),
		ProjectID:     projectID,
		Name:          req.Name,
		Priority:      req.Priority,
		StartDateTime: req.StartDateTime,
		EndDateTime:   req.EndDateTime,
		Status:        "todo",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if req.Description != nil {
		task.Description = req.Description
	}

	if req.Location != nil {
		task.Location = req.Location
	}

	if req.RecurringDays != nil {
		task.RecurringDays = req.RecurringDays
	}

	if req.RecurringUntil != nil {
		task.RecurringUntil = req.RecurringUntil
	}

	// persist account to database
	err := s.repo.CreateTask(ctx, task)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to create task", "CREATE_TASK_ERROR", err)
	}

	return task, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID uuid.UUID) (*entity.Task, error) {
	tsk, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("task not found", "TASK_NOT_FOUND", err)
		}
		return nil, apperrors.NewInternalServerError("failed to get task", "GET_TASK_ERROR", err)
	}

	return tsk, nil
}

func (s *TaskService) ListTasksByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Task, error) {
	tasks, err := s.repo.ListTasksByProject(ctx, projectID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to list tasks", "LIST_TASKS_ERROR", err)
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID uuid.UUID, req *task.UpdateTaskRequest) (*entity.Task, error) {
	// Get task by ID for update
	tsk, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("task not found", "TASK_NOT_FOUND", err)
		}
		return nil, apperrors.NewInternalServerError("failed to get task", "GET_TASK_ERROR", err)
	}

	// Rule: If status != todo, cannot change start_datetime
	if tsk.Status != "todo" && req.StartDateTime != nil {
		return nil, apperrors.NewBadRequestError("cannot update start_datetime when status is not todo", "INVALID_REQUEST", nil)
	}

	// Additional validation same as CreateTask
	if req.StartDateTime != nil && req.EndDateTime != nil {
		if *req.StartDateTime == *req.EndDateTime {
			return nil, apperrors.NewBadRequestError("start_datetime and end_datetime cannot be the same", "INVALID_REQUEST", nil)
		}
		if *req.EndDateTime < *req.StartDateTime {
			return nil, apperrors.NewBadRequestError("end_datetime must be greater than start_datetime", "INVALID_REQUEST", nil)
		}
	}

	// Update fields
	tsk.Name = req.Name
	tsk.Description = req.Description
	tsk.Priority = req.Priority
	tsk.Location = req.Location
	tsk.RecurringDays = req.RecurringDays
	tsk.RecurringUntil = req.RecurringUntil
	tsk.StartDateTime = req.StartDateTime
	tsk.EndDateTime = req.EndDateTime
	tsk.UpdatedAt = time.Now()

	err = s.repo.UpdateTask(ctx, tsk)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to update task", "UPDATE_TASK_ERROR", err)
	}

	return tsk, nil
}
