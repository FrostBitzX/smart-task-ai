package service

import (
	"context"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
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
