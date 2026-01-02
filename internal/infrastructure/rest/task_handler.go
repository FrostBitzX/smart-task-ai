package rest

import (
	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	"github.com/FrostBitzX/smart-task-ai/internal/application/task/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	CreateTaskUC *usecase.CreateTaskUseCase
	logger       logger.Logger
}

func NewTaskHandler(
	create *usecase.CreateTaskUseCase,
	l logger.Logger,
) *TaskHandler {
	return &TaskHandler{
		CreateTaskUC: create,
		logger:       l,
	}
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	req, err := requests.ParseAndValidate[task.CreateTaskRequest](c)
	if err != nil {
		h.logger.Warn("Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return responses.Error(c, apperror.ErrInvalidData)
	}

	// Parse project ID from URL
	projectID := c.Params("projectID")
	if projectID == "" {
		return responses.Error(c, apperrors.NewBadRequestError("project ID is required", "INVALID_PROJECT_ID", nil))
	}

	data, err := h.CreateTaskUC.Execute(c.Context(), projectID, req)
	if err != nil {
		return responses.Error(c, err)
	}

	return responses.Success(c, data, "Task created successfully")
}
