package routes

import (
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/groq"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/middlewares"

	chatUC "github.com/FrostBitzX/smart-task-ai/internal/application/chat/usecase"
	profileUC "github.com/FrostBitzX/smart-task-ai/internal/application/profile/usecase"
	projectUC "github.com/FrostBitzX/smart-task-ai/internal/application/project/usecase"
	taskUC "github.com/FrostBitzX/smart-task-ai/internal/application/task/usecase"
	chatDomain "github.com/FrostBitzX/smart-task-ai/internal/domain/chats/service"
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
	taskRepository := repo.NewTaskRepository(db)
	projectService := projectDomain.NewProjectService(projectRepository, taskRepository)
	createProjectUC := projectUC.NewCreateProjectUseCase(projectService, log)
	listProjectByAccountUC := projectUC.NewListProjectByAccountUseCase(projectService, log)
	getProjectByIDUC := projectUC.NewGetProjectByIDUseCase(projectService, log)
	updateProjectUC := projectUC.NewUpdateProjectUseCase(projectService, log)
	deleteProjectUC := projectUC.NewDeleteProjectUseCase(projectService, log)
	projectHandlerInstance := handler.NewProjectHandler(
		createProjectUC,
		listProjectByAccountUC,
		getProjectByIDUC,
		updateProjectUC,
		deleteProjectUC,
		log,
	)

	// Project routes
	api.Post("/projects", projectHandlerInstance.CreateProject)
	api.Get("/projects", projectHandlerInstance.ListProject)
	api.Get("/projects/:projectId", projectHandlerInstance.GetProject)
	api.Patch("/projects/:projectId", projectHandlerInstance.UpdateProject)
	api.Delete("/projects/:projectId", projectHandlerInstance.DeleteProject)

	// Task setup
	taskService := taskDomain.NewTaskService(taskRepository)
	createTaskUC := taskUC.NewCreateTaskUseCase(taskService, log)
	getTaskByIDUC := taskUC.NewGetTaskByIDUseCase(taskService, log)
	listTasksByProjectUC := taskUC.NewListTasksByProjectUseCase(taskService, log)
	updateTaskUC := taskUC.NewUpdateTaskUseCase(taskService, log)
	deleteTaskUC := taskUC.NewDeleteTaskUseCase(taskService, log)
	taskHandlerInstance := handler.NewTaskHandler(createTaskUC, getTaskByIDUC, listTasksByProjectUC, updateTaskUC, deleteTaskUC, log)

	// Task routes
	api.Post("/:projectId/tasks", taskHandlerInstance.CreateTask)
	api.Get("/:projectId/tasks", taskHandlerInstance.ListTasksByProject)
	api.Get("/tasks/:taskId", taskHandlerInstance.GetTaskByID)
	api.Patch("/tasks/:taskId", taskHandlerInstance.UpdateTask)
	api.Delete("/tasks/:taskId", taskHandlerInstance.DeleteTask)

	// Chat setup
	groqClient, err := groq.NewGroqClient()
	if err != nil {
		log.Warn("Failed to initialize Groq client, chat endpoints will not be available", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		chatService := chatDomain.NewChatService(groqClient, taskService, projectService)
		sendMessageUC := chatUC.NewSendMessageUseCase(chatService, log)
		chatHandlerInstance := handler.NewChatHandler(sendMessageUC, log)

		// Chat routes (protected by JWT middleware via /api group)
		api.Post("/:projectId/chat", chatHandlerInstance.SendMessage)
	}
}
