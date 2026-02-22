package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/tasks/service"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/mocks"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// mockLogger is a simple logger implementation for testing
type mockLogger struct{}

func (m *mockLogger) Info(msg string, fields ...map[string]interface{})  {}
func (m *mockLogger) Error(msg string, fields ...map[string]interface{}) {}
func (m *mockLogger) Debug(msg string, fields ...map[string]interface{}) {}
func (m *mockLogger) Warn(msg string, fields ...map[string]interface{})  {}
func (m *mockLogger) With(fields map[string]interface{}) logger.Logger   { return m }

func TestGetTaskStatisticsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	mockProjectRepo := mocks.NewMockProjectRepository(ctrl)
	taskService := service.NewTaskService(mockRepo, mockProjectRepo)
	logger := &mockLogger{}
	uc := NewGetTaskStatisticsUseCase(taskService, logger)
	ctx := context.Background()
	nodeIDStr := "550e8400-e29b-41d4-a716-446655440000"
	nodeID := uuid.MustParse(nodeIDStr)

	tests := []struct {
		name          string
		nodeIDStr     string
		setupMock     func()
		expectedStats *struct {
			Todo       int64
			InProgress int64
			InReview   int64
			Done       int64
		}
		expectedError string
	}{
		{
			name:      "success - returns statistics with all statuses",
			nodeIDStr: nodeIDStr,
			setupMock: func() {
				statusCounts := []tasks.StatusCount{
					{Status: "todo", Count: 5},
					{Status: "in_progress", Count: 3},
					{Status: "in_review", Count: 2},
					{Status: "done", Count: 10},
				}
				mockRepo.EXPECT().
					CountTasksByStatus(ctx, nodeID).
					Return(statusCounts, nil).
					Times(1)
			},
			expectedStats: &struct {
				Todo       int64
				InProgress int64
				InReview   int64
				Done       int64
			}{
				Todo:       5,
				InProgress: 3,
				InReview:   2,
				Done:       10,
			},
			expectedError: "",
		},
		{
			name:      "success - returns zero statistics when no tasks",
			nodeIDStr: nodeIDStr,
			setupMock: func() {
				mockRepo.EXPECT().
					CountTasksByStatus(ctx, nodeID).
					Return([]tasks.StatusCount{}, nil).
					Times(1)
			},
			expectedStats: &struct {
				Todo       int64
				InProgress int64
				InReview   int64
				Done       int64
			}{
				Todo:       0,
				InProgress: 0,
				InReview:   0,
				Done:       0,
			},
			expectedError: "",
		},
		{
			name:      "error - invalid node ID format",
			nodeIDStr: "invalid-uuid",
			setupMock: func() {
				// No mock setup needed as validation happens before repository call
			},
			expectedStats: nil,
			expectedError: "invalid node ID format",
		},
		{
			name:      "error - repository fails",
			nodeIDStr: nodeIDStr,
			setupMock: func() {
				mockRepo.EXPECT().
					CountTasksByStatus(ctx, nodeID).
					Return(nil, apperror.NewInternalServerError("database error", "DB_ERROR", errors.New("connection failed"))).
					Times(1)
			},
			expectedStats: nil,
			expectedError: "failed to get task statistics",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			res, err := uc.Execute(ctx, tt.nodeIDStr)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, res)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tt.expectedStats.Todo, res.Todo)
				assert.Equal(t, tt.expectedStats.InProgress, res.InProgress)
				assert.Equal(t, tt.expectedStats.InReview, res.InReview)
				assert.Equal(t, tt.expectedStats.Done, res.Done)
			}
		})
	}
}
