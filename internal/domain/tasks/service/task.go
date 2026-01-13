package service

import (
	"context"
	"errors"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"

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
		return nil, apperror.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	if err := s.validateTimeRange(req.StartDateTime, req.EndDateTime); err != nil {
		return nil, err
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
		return nil, apperror.NewInternalServerError("failed to create task", "CREATE_TASK_ERROR", err)
	}

	return task, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID uuid.UUID) (*entity.Task, error) {
	tsk, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("task not found", "TASK_NOT_FOUND", err)
		}
		return nil, apperror.NewInternalServerError("failed to get task", "GET_TASK_ERROR", err)
	}

	return tsk, nil
}

func (s *TaskService) ListTasksByProject(ctx context.Context, projectID uuid.UUID) ([]*entity.Task, error) {
	tasks, err := s.repo.ListTasksByProject(ctx, projectID)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to list tasks", "LIST_TASKS_ERROR", err)
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID uuid.UUID, req *task.UpdateTaskRequest) (*entity.Task, error) {
	// Get task by ID for update
	tsk, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("task not found", "TASK_NOT_FOUND", err)
		}
		return nil, apperror.NewInternalServerError("failed to get task", "GET_TASK_ERROR", err)
	}

	// Rule: If status != todo, cannot change start_datetime
	if tsk.Status != "todo" && req.StartDateTime != nil {
		return nil, apperror.NewBadRequestError("cannot update start_datetime when status is not todo", "INVALID_REQUEST", nil)
	}

	// Additional validation same as CreateTask
	if err := s.validateTimeRange(req.StartDateTime, req.EndDateTime); err != nil {
		return nil, err
	}

	// Update fields only if provided (PATCH semantics)
	if req.Name != "" {
		tsk.Name = req.Name
	}
	if req.Status != nil {
		tsk.Status = *req.Status
	}
	if req.Description != nil {
		tsk.Description = req.Description
	}
	if req.Priority != "" {
		tsk.Priority = req.Priority
	}
	if req.Location != nil {
		tsk.Location = req.Location
	}
	if req.RecurringDays != nil {
		tsk.RecurringDays = req.RecurringDays
	}
	if req.RecurringUntil != nil {
		tsk.RecurringUntil = req.RecurringUntil
	}
	if req.StartDateTime != nil {
		tsk.StartDateTime = req.StartDateTime
	}
	if req.EndDateTime != nil {
		tsk.EndDateTime = req.EndDateTime
	}
	tsk.UpdatedAt = time.Now()

	err = s.repo.UpdateTask(ctx, tsk)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to update task", "UPDATE_TASK_ERROR", err)
	}

	return tsk, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	_, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return apperror.NewNotFoundError("task not found", "TASK_NOT_FOUND", err)
		}
		return apperror.NewInternalServerError("failed to get task", "GET_TASK_ERROR", err)
	}

	err = s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		return apperror.NewInternalServerError("failed to delete task", "DELETE_TASK_ERROR", err)
	}

	return nil
}

func (s *TaskService) validateTimeRange(startStr, endStr *string) error {
	if startStr != nil && endStr != nil {
		start, err := time.Parse(time.RFC3339, *startStr)
		if err != nil {
			return apperror.NewBadRequestError("invalid start_datetime format", "INVALID_DATE_FORMAT", err)
		}
		end, err := time.Parse(time.RFC3339, *endStr)
		if err != nil {
			return apperror.NewBadRequestError("invalid end_datetime format", "INVALID_DATE_FORMAT", err)
		}

		if start.Equal(end) {
			return apperror.NewBadRequestError("start_datetime and end_datetime cannot be the same", "INVALID_REQUEST", nil)
		}
		if end.Before(start) {
			return apperror.NewBadRequestError("end_datetime must be greater than start_datetime", "INVALID_REQUEST", nil)
		}
	}
	return nil
}
