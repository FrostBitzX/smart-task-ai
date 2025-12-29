package service

import (
	"context"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/project"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
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
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// create domain entity
	now := time.Now()
	proj := &entity.Project{
		ID:        uuid.New(),
		AccountID: uuid.MustParse(req.AccountID),
		Role:      "owner",
		Name:      req.Name,
		Config:    req.Config,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// persist account to database
	err := s.repo.CreateProject(ctx, proj)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to create project", "CREATE_PROJECT_ERROR", err)
	}

	return proj, nil
}
