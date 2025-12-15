package routes

import (
	accSignUpUC "github.com/FrostBitzX/smart-task-ai/internal/application/account/usecase"
	accDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/service"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	repo "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/persistence"
	accHandler "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/rest"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/middlewares"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPublicRoutes(app fiber.Router, db *gorm.DB, log logger.Logger) {
	api := app.Group("/api")

	accountRepository := repo.NewAccountRepository(db)
	accountService := accDomain.NewAccountService(accountRepository)
	accountSignUpUC := accSignUpUC.NewAccountUseCase(accountService, log)
	accountHandler := accHandler.NewAccountHandler(accountSignUpUC, log)

	api.Post("/account", middlewares.ValidateCreateAccountRequest, accountHandler.CreateAccount)
}
