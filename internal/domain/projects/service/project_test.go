package service

import (
	"context"
	"testing"

	"github.com/FrostBitzX/smart-task-ai/internal/application/project"
	"github.com/FrostBitzX/smart-task-ai/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProjectService_CreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProjectRepository(ctrl)
	svc := NewProjectService(mockRepo)

	ctx := context.Background()
	accountID := "550e8400-e29b-41d4-a716-446655440000" // Valid UUID
	req := &project.CreateProjectRequest{
		AccountID: accountID,
		Name:      "Test Project",
		Config:    []byte(`{"color": "blue"}`),
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			CreateProject(ctx, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := svc.CreateProject(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Name, res.Name)
		assert.Equal(t, req.Config, res.Config)
	})

	t.Run("invalid request body", func(t *testing.T) {
		res, err := svc.CreateProject(ctx, nil)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid request body")
	})

	t.Run("invalid account ID format", func(t *testing.T) {
		invalidReq := &project.CreateProjectRequest{
			AccountID: "invalid-uuid",
		}
		res, err := svc.CreateProject(ctx, invalidReq)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "invalid account ID format")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.EXPECT().
			CreateProject(ctx, gomock.Any()).
			Return(assert.AnError).
			Times(1)

		res, err := svc.CreateProject(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
