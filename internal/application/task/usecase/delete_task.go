package usecase

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/service"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
)

type DeleteTaskUseCase struct {
	taskService *service.TaskService
	logger      logger.Logger
}

func NewDeleteTaskUseCase(s *service.TaskService, l logger.Logger) *DeleteTaskUseCase {
	return &DeleteTaskUseCase{
		taskService: s,
		logger:      l,
	}
}

func (uc *DeleteTaskUseCase) Execute(ctx context.Context, taskID string) (string, error) {
	parsedTaskID, err := utils.ParseID(taskID, entity.TaskIDPrefix)
	if err != nil {
		return "", apperrors.NewBadRequestError("invalid task ID format", "INVALID_TASK_ID", err)
	}

	err = uc.taskService.DeleteTask(ctx, parsedTaskID)
	if err != nil {
		return "", err
	}

	return taskID, nil
}
