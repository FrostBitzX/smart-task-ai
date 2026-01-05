package rest

import (
	"github.com/FrostBitzX/smart-task-ai/internal/application/project"
	"github.com/FrostBitzX/smart-task-ai/internal/application/project/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	CreateProjectUC *usecase.CreateProjectUseCase
	logger          logger.Logger
}

func NewProjectHandler(
	create *usecase.CreateProjectUseCase,
	l logger.Logger,
) *ProjectHandler {
	return &ProjectHandler{
		CreateProjectUC: create,
		logger:          l,
	}
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	req, err := requests.ParseAndValidate[project.CreateProjectRequest](c)
	if err != nil {
		h.logger.Warn("Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return responses.Error(c, err)
	}

	// Get AccountID from JWT claims
	jwtClaims, ok := c.Locals("jwt_claims").(map[string]interface{})
	if !ok {
		h.logger.Error("Invalid JWT claims", nil)
		return responses.Error(c, apperror.ErrUnauthorized)
	}

	accountID, ok := jwtClaims["AccountId"].(string)
	if !ok || accountID == "" {
		h.logger.Error("Missing AccountId in JWT claims", nil)
		return responses.Error(c, apperror.ErrUnauthorized)
	}

	// Set AccountID from JWT
	req.AccountID = accountID

	data, err := h.CreateProjectUC.Execute(c.Context(), req)
	if err != nil {
		return responses.Error(c, err)
	}

	return responses.Success(c, data, "Project created successfully")
}
