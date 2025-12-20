package rest

import (
	"github.com/FrostBitzX/smart-task-ai/internal/application/profile"
	"github.com/FrostBitzX/smart-task-ai/internal/application/profile/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

type ProfileHandler struct {
	CreateProfileUC *usecase.CreateProfileUseCase
	logger          logger.Logger
}

func NewProfileHandler(
	create *usecase.CreateProfileUseCase,
	l logger.Logger,
) *ProfileHandler {
	return &ProfileHandler{
		CreateProfileUC: create,
		logger:          l,
	}

}

func (h *ProfileHandler) CreateProfile(c *fiber.Ctx) error {
	// Get AccountID from JWT claims
	jwtClaims := c.Locals("jwt_claims").(map[string]interface{})
	accountID := jwtClaims["AccountId"].(string)

	req, err := requests.ParseAndValidate[profile.CreateProfileRequest](c)
	if err != nil {
		h.logger.Warn("Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return responses.Error(c, apperror.ErrInvalidData)
	}

	// Set AccountID from JWT
	req.AccountID = accountID

	data, err := h.CreateProfileUC.Execute(c.Context(), req)
	if err != nil {
		return responses.Error(c, err)
	}

	return responses.Success(c, data, "Profile created successfully")
}
