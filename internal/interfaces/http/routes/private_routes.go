package routes

import (
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/middlewares"

	profileUC "github.com/FrostBitzX/smart-task-ai/internal/application/profile/usecase"
	profileDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/service"
	repo "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/persistence"
	profileHandler "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/rest"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterPrivateRoutes(app fiber.Router, db *gorm.DB, log logger.Logger) {
	api := app.Group("/api", middlewares.JWTMiddleware())

	// Profile setup
	profileRepository := repo.NewProfileRepository(db)
	profileService := profileDomain.NewProfileService(profileRepository)
	createProfileUC := profileUC.NewCreateProfileUseCase(profileService, log)
	getProfileUC := profileUC.NewGetProfileUseCase(profileService, log)
	profileHandlerInstance := profileHandler.NewProfileHandler(createProfileUC, getProfileUC, log)

	// Profile routes
	api.Post("/profile", profileHandlerInstance.CreateProfile)
	api.Get("/profile", profileHandlerInstance.GetProfile)
}
