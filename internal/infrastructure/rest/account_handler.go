package account

import (
	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/application/account/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	AccountUC *usecase.AccountUseCase
	logger    logger.Logger
}

func NewAccountHandler(accUC *usecase.AccountUseCase, l logger.Logger) *AccountHandler {
	return &AccountHandler{
		AccountUC: accUC,
		logger:    l,
	}
}

func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	req, err := requests.ParseAndValidate[account.CreateAccountRequest](c)
	if err != nil {
		h.logger.Warn("Invalid request data", map[string]interface{}{
			"error": err.Error(),
		})
		return responses.Error(c, apperror.ErrInvalidData)
	}

	data, err := h.AccountUC.Execute(req)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(data)
}
