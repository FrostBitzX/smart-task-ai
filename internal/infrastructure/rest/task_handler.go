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
	CreateTaskUC  *usecase.CreateTaskUseCase
	GetTaskByIDUC *usecase.GetTaskByIDUseCase
	logger        logger.Logger
}

func NewTaskHandler(
	create *usecase.CreateTaskUseCase,
	getByID *usecase.GetTaskByIDUseCase,
	l logger.Logger,
) *TaskHandler {
	return &TaskHandler{
		CreateTaskUC:  create,
		GetTaskByIDUC: getByID,
		logger:        l,
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

func (h *TaskHandler) GetTaskByID(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	if taskID == "" {
		return responses.Error(c, apperrors.NewBadRequestError("task ID is required", "INVALID_TASK_ID", nil))
	}

	data, err := h.GetTaskByIDUC.Execute(c.Context(), taskID)
	if err != nil {
		return responses.Error(c, err)
	}

	return responses.Success(c, data, "Task retrieved successfully")
}
