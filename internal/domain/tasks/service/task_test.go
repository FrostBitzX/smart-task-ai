package service

import (
	"context"
	"testing"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/task"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTaskService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	svc := NewTaskService(mockRepo)
	ctx := context.Background()
	projectID := uuid.New()

	t.Run("success", func(t *testing.T) {
		req := &task.CreateTaskRequest{
			Name:     "Test Task",
			Priority: "1",
		}
		mockRepo.EXPECT().
			CreateTask(ctx, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := svc.CreateTask(ctx, projectID, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.Name, res.Name)
	})

	t.Run("invalid time range - same time", func(t *testing.T) {
		nowStr := time.Now().Format(time.RFC3339)
		req := &task.CreateTaskRequest{
			Name:          "Test Task",
			StartDateTime: &nowStr,
			EndDateTime:   &nowStr,
		}

		res, err := svc.CreateTask(ctx, projectID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "cannot be the same")
	})

	t.Run("repository error", func(t *testing.T) {
		req := &task.CreateTaskRequest{
			Name: "Test Task",
		}
		mockRepo.EXPECT().
			CreateTask(ctx, gomock.Any()).
			Return(assert.AnError).
			Times(1)

		res, err := svc.CreateTask(ctx, projectID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskService_GetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	svc := NewTaskService(mockRepo)
	ctx := context.Background()
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		expectedTask := &entity.Task{ID: taskID, Name: "Test Task"}
		mockRepo.EXPECT().
			GetTaskByID(ctx, taskID).
			Return(expectedTask, nil).
			Times(1)

		res, err := svc.GetTaskByID(ctx, taskID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTask, res)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.EXPECT().
			GetTaskByID(ctx, taskID).
			Return(nil, assert.AnError).
			Times(1)

		res, err := svc.GetTaskByID(ctx, taskID)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestTaskService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	svc := NewTaskService(mockRepo)
	ctx := context.Background()
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		existingTask := &entity.Task{ID: taskID, Status: "todo", Name: "Old Name"}
		req := &task.UpdateTaskRequest{Name: "New Name"}

		mockRepo.EXPECT().
			GetTaskByID(ctx, taskID).
			Return(existingTask, nil).
			Times(1)
		mockRepo.EXPECT().
			UpdateTask(ctx, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := svc.UpdateTask(ctx, taskID, req)

		assert.NoError(t, err)
		assert.Equal(t, "New Name", res.Name)
	})

	t.Run("cannot update start_datetime when not todo", func(t *testing.T) {
		existingTask := &entity.Task{ID: taskID, Status: "doing"}
		newDate := time.Now().Format(time.RFC3339)
		req := &task.UpdateTaskRequest{StartDateTime: &newDate}

		mockRepo.EXPECT().
			GetTaskByID(ctx, taskID).
			Return(existingTask, nil).
			Times(1)

		res, err := svc.UpdateTask(ctx, taskID, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "cannot update start_datetime")
	})
}

func TestTaskService_DeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	svc := NewTaskService(mockRepo)
	ctx := context.Background()
	taskID := uuid.New()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetTaskByID(ctx, taskID).
			Return(&entity.Task{}, nil).
			Times(1)
		mockRepo.EXPECT().
			DeleteTask(ctx, taskID).
			Return(nil).
			Times(1)

		err := svc.DeleteTask(ctx, taskID)

		assert.NoError(t, err)
	})
}
