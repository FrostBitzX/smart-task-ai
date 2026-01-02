package routes

import (
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/middlewares"

	profileUC "github.com/FrostBitzX/smart-task-ai/internal/application/profile/usecase"
	projectUC "github.com/FrostBitzX/smart-task-ai/internal/application/project/usecase"
	taskUC "github.com/FrostBitzX/smart-task-ai/internal/application/task/usecase"
	profileDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/service"
	projectDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/projects/service"
	taskDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/service"
	repo "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/persistence"
	handler "github.com/FrostBitzX/smart-task-ai/internal/infrastructure/rest"

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
	updateProfileUC := profileUC.NewUpdateProfileUseCase(profileService, log)
	profileHandlerInstance := handler.NewProfileHandler(createProfileUC, getProfileUC, updateProfileUC, log)

	// Profile routes
	api.Post("/profiles", profileHandlerInstance.CreateProfile)
	api.Get("/profiles", profileHandlerInstance.GetProfile)
	api.Patch("/profiles", profileHandlerInstance.UpdateProfile)

	// Project setup
	projectRepository := repo.NewProjectRepository(db)
	projectService := projectDomain.NewProjectService(projectRepository)
	createProjectUC := projectUC.NewCreateProjectUseCase(projectService, log)
	projectHandlerInstance := handler.NewProjectHandler(createProjectUC, log)

	// Project routes
	api.Post("/projects", projectHandlerInstance.CreateProject)

	// Task setup
	taskRepository := repo.NewTaskRepository(db)
	taskService := taskDomain.NewTaskService(taskRepository)
	createTaskUC := taskUC.NewCreateTaskUseCase(taskService, log)
	getTaskByIDUC := taskUC.NewGetTaskByIDUseCase(taskService, log)
	taskHandlerInstance := handler.NewTaskHandler(createTaskUC, getTaskByIDUC, log)

	// Task routes
	api.Post("/:projectId/tasks", taskHandlerInstance.CreateTask)
	api.Get("/tasks/:taskId", taskHandlerInstance.GetTaskByID)
}
