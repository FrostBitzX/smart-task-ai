package routes

import (
	accUseCase "github.com/FrostBitzX/smart-task-ai/internal/application/account/usecase"
	accDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/accounts/service"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	repo "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/persistence"
	accHandler "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/rest"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPublicRoutes(app fiber.Router, db *gorm.DB, log logger.Logger) {
	api := app.Group("/api")

	accountRepository := repo.NewAccountRepository(db)
	accountService := accDomain.NewAccountService(accountRepository)
	accountUsecase := accUseCase.NewAccountUseCase(accountService, log)
	accountHandler := accHandler.NewAccountHandler(accountUsecase, log)

	api.Post("/signup", accountHandler.CreateAccount)
}
