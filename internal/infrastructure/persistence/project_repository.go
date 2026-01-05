package persistence

import (
	"context"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/projects/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) projects.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(ctx context.Context, proj *entity.Project) error {
	return r.db.WithContext(ctx).Create(proj).Error
}

func (r *projectRepository) GetProjectByID(ctx context.Context, projectID uuid.UUID) (*entity.Project, error) {
	var proj entity.Project
	err := r.db.WithContext(ctx).
		Where("id = ?", projectID).
		First(&proj).Error
	if err != nil {
		return nil, err
	}
	return &proj, nil
}
