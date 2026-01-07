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

func (r *projectRepository) ListProjectByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]*entity.Project, int, error) {
	var projects []*entity.Project
	var total int64

	// Get total count
	if err := r.db.WithContext(ctx).Model(&entity.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, int(total), err
}
