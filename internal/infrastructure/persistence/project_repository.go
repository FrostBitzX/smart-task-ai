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

func (r *projectRepository) GetProjectByID(ctx context.Context, projectID uuid.UUID, nodeID uuid.UUID) (*entity.Project, error) {
	var proj entity.Project
	err := r.db.WithContext(ctx).
		Where("id = ? AND node_id = ?", projectID, nodeID).
		First(&proj).Error
	if err != nil {
		return nil, err
	}
	return &proj, nil
}

func (r *projectRepository) ListProjectByAccountID(ctx context.Context, accountID uuid.UUID, nodeID uuid.UUID, limit, offset int) ([]*entity.Project, int, error) {
	var projects []*entity.Project
	var total int64

	// Get total count with node_id filter
	if err := r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Where("account_id = ? AND node_id = ?", accountID, nodeID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with node_id filter
	err := r.db.WithContext(ctx).
		Where("account_id = ? AND node_id = ?", accountID, nodeID).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, int(total), err
}

func (r *projectRepository) UpdateProject(ctx context.Context, proj *entity.Project, nodeID uuid.UUID) error {
	// Verify the project belongs to the tenant before updating
	return r.db.WithContext(ctx).
		Model(&entity.Project{}).
		Select("*").
		Where("id = ? AND node_id = ?", proj.ID, nodeID).
		Updates(proj).Error
}

func (r *projectRepository) DeleteProject(ctx context.Context, projectID uuid.UUID, nodeID uuid.UUID) error {
	// Verify the project belongs to the tenant before deleting
	return r.db.WithContext(ctx).
		Where("id = ? AND node_id = ?", projectID, nodeID).
		Delete(&entity.Project{}).Error
}

func (r *projectRepository) AddMember(ctx context.Context, member *entity.ProjectMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *projectRepository) RemoveMember(ctx context.Context, projectID uuid.UUID, accountID uuid.UUID, nodeID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("project_id = ? AND account_id = ? AND node_id = ? AND role != ?", projectID, accountID, nodeID, "owner").
		Delete(&entity.ProjectMember{}).Error
}

func (r *projectRepository) GetProjectMember(ctx context.Context, projectID uuid.UUID, accountID uuid.UUID, nodeID uuid.UUID) (*entity.ProjectMember, error) {
	var member entity.ProjectMember
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND account_id = ? AND node_id = ?", projectID, accountID, nodeID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *projectRepository) ListProjectMembers(ctx context.Context, projectID uuid.UUID, nodeID uuid.UUID) ([]*entity.ProjectMember, error) {
	var members []*entity.ProjectMember
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND node_id = ?", projectID, nodeID).
		Find(&members).Error
	return members, err
}
