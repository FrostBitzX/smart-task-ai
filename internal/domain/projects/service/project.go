package service

import (
	"context"
	"errors"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/project"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
)

type ProjectService struct {
	repo projects.ProjectRepository
}

func NewProjectService(repo projects.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
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
