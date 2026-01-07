package service

import (
	"context"
	"errors"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/project"
	"github.com/FrostBitzX/smart-task-ai/internal/utils"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks"
)

type ProjectService struct {
	repo     projects.ProjectRepository
	taskRepo tasks.TaskRepository
}

func NewProjectService(repo projects.ProjectRepository, taskRepo tasks.TaskRepository) *ProjectService {
	return &ProjectService{
		repo:     repo,
		taskRepo: taskRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *project.CreateProjectRequest) (*entity.Project, error) {
	if req == nil {
		return nil, apperror.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// create domain entity
	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid account ID format", "INVALID_ACCOUNT_ID", err)
	}

	now := time.Now()
	proj := &entity.Project{
		ID:        uuid.New(),
		AccountID: accountID,
		Role:      "owner",
		Name:      req.Name,
		Config:    req.Config,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// persist account to database
	err = s.repo.CreateProject(ctx, proj)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to create project", "CREATE_PROJECT_ERROR", err)
	}

	return proj, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, projectID uuid.UUID) (*entity.Project, error) {
	proj, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("project not found", "PROJECT_NOT_FOUND", err)
		}
		return nil, apperror.NewInternalServerError("failed to get project", "GET_PROJECT_ERROR", err)
	}

	return proj, nil
}

func (s *ProjectService) ListProjectByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Project, int, error) {
	projs, total, err := s.repo.ListProjectByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("failed to list projects", "LIST_PROJECT_ERROR", err)
	}

	return projs, total, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, req *project.UpdateProjectRequest) (*entity.Project, error) {
	if req == nil {
		return nil, apperror.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	projectID, err := utils.ParseID(req.ProjectID, entity.ProjectIDPrefix)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid project ID format", "INVALID_PROJECT_ID", err)
	}

	// Get existing project
	proj, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("project not found", "PROJECT_NOT_FOUND", err)
		}
		return nil, apperror.NewInternalServerError("failed to get project", "GET_PROJECT_ERROR", err)
	}

	proj.Name = req.Name
	if req.Config != nil {
		proj.Config = req.Config
	}

	proj.UpdatedAt = time.Now()

	err = s.repo.UpdateProject(ctx, proj)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to update project", "UPDATE_PROJECT_ERROR", err)
	}

	return proj, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	// Get existing project
	_, err := s.repo.GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return apperror.NewNotFoundError("project not found", "PROJECT_NOT_FOUND", err)
		}
		return apperror.NewInternalServerError("failed to get project", "GET_PROJECT_ERROR", err)
	}

	if err := s.deleteProjectCheck(ctx, projectID); err != nil {
		return err
	}

	err = s.repo.DeleteProject(ctx, projectID)
	if err != nil {
		return apperror.NewInternalServerError("failed to delete project", "DELETE_PROJECT_ERROR", err)
	}

	return nil
}

func (s *ProjectService) deleteProjectCheck(ctx context.Context, projectID uuid.UUID) error {
	count, err := s.taskRepo.CountTasksByProject(ctx, projectID)
	if err != nil {
		return apperror.NewInternalServerError("failed to check tasks in project", "CHECK_TASKS_ERROR", err)
	}

	if count > 0 {
		return apperror.NewBadRequestError("cannot delete project with existing tasks", "PROJECT_HAS_TASKS", nil)
	}

	return nil
}
