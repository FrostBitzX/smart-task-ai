package routes

import (
	accUC "github.com/FrostBitzX/smart-task-ai/internal/application/account/usecase"
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
	accountSignUpUC := accUC.NewAccountUseCase(accountService, log)
	accountLoginUC := accUC.NewLoginUseCase(accountService, log)
	accountHandler := accHandler.NewAccountHandler(accountSignUpUC, accountLoginUC, log)

	api.Post("/signup", accountHandler.CreateAccount)
	api.Post("/login", accountHandler.Login)
}
