package account

import (
	"github.com/FrostBitzX/smart-task-ai/internal/application/account"
	"github.com/FrostBitzX/smart-task-ai/internal/application/account/usecase"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/middlewares"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/requests"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	AccountUC     *usecase.AccountUseCase
	ListAccountUC *usecase.ListAccountUseCase
	logger        logger.Logger
}

func NewAccountHandler(accUC *usecase.AccountUseCase, listUC *usecase.ListAccountUseCase, l logger.Logger) *AccountHandler {
	return &AccountHandler{
		AccountUC:     accUC,
		ListAccountUC: listUC,
		logger:        l,
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
		return responses.Error(c, err)
	}

	return responses.Success(c, data, "Account created successfully")
}

func (h *AccountHandler) ListAccounts(c *fiber.Ctx) error {
	// Get validated request from middleware context
	req, ok := c.Locals("validatedRequest").(middlewares.ListAccountsRequest)
	if !ok {
		h.logger.Warn("Failed to get validated request from context")
		return responses.Error(c, apperror.ErrInvalidData)
	}

	// Convert middleware request to domain request
	domainReq := &account.ListAccountsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	data, err := h.ListAccountUC.Execute(domainReq)
	if err != nil {
		return responses.Error(c, err)
	}

	return responses.Success(c, data, "List accounts successfully")
}
