package usecase

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/service"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
)

type UpdateTaskUseCase struct {
	taskService *service.TaskService
	logger      logger.Logger
}

func NewUpdateTaskUseCase(s *service.TaskService, l logger.Logger) *UpdateTaskUseCase {
	return &UpdateTaskUseCase{
		taskService: s,
		logger:      l,
	}
}

func (uc *UpdateTaskUseCase) Execute(ctx context.Context, taskID string, req *task.UpdateTaskRequest) (*task.UpdateTaskResponse, error) {
	if req == nil {
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	parsedTaskID, err := utils.ParseID(taskID, entity.TaskIDPrefix)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid task ID format", "INVALID_TASK_ID", err)
	}

	result, err := uc.taskService.UpdateTask(ctx, parsedTaskID, req)
	if err != nil {
		return nil, err
	}

	// Convert UUID to string with prefix
	taskIDRes := utils.ShortUUIDWithPrefix(result.ID, entity.TaskIDPrefix)

	return &task.UpdateTaskResponse{
		ID:             taskIDRes,
		Status:         result.Status,
		Name:           result.Name,
		Description:    result.Description,
		Priority:       result.Priority,
		StartDateTime:  result.StartDateTime,
		EndDateTime:    result.EndDateTime,
		Location:       result.Location,
		RecurringDays:  result.RecurringDays,
		RecurringUntil: result.RecurringUntil,
	}, nil
}
