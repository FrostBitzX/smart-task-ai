package service

import (
	"context"
	"testing"

	"github.com/FrostBitzX/smart-task-ai/internal/application/profile"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
	"github.com/FrostBitzX/smart-task-ai/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProfileService_CreateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProfileRepository(ctrl)
	svc := NewProfileService(mockRepo)
	ctx := context.Background()
	accountID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		req := &profile.CreateProfileRequest{
			AccountID: accountID,
			FirstName: "First",
			LastName:  "Last",
		}

		mockRepo.EXPECT().
			GetProfileByAccountID(ctx, accountID).
			Return(nil, nil).
			Times(1)
		mockRepo.EXPECT().
			CreateProfile(ctx, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := svc.CreateProfile(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.FirstName, res.FirstName)
	})

	t.Run("profile already exists", func(t *testing.T) {
		req := &profile.CreateProfileRequest{
			AccountID: accountID,
		}

		mockRepo.EXPECT().
			GetProfileByAccountID(ctx, accountID).
			Return(&entity.Profile{}, nil).
			Times(1)

		res, err := svc.CreateProfile(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "profile already exists")
	})
}

func TestProfileService_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockProfileRepository(ctrl)
	svc := NewProfileService(mockRepo)
	ctx := context.Background()
	accountID := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		existingProfile := &entity.Profile{ID: uuid.New(), AccountID: uuid.MustParse(accountID)}
		req := &profile.UpdateProfileRequest{
			AccountID: accountID,
			FirstName: "Updated Name",
		}

		mockRepo.EXPECT().
			GetProfileByAccountID(ctx, accountID).
			Return(existingProfile, nil).
			Times(1)
		mockRepo.EXPECT().
			UpdateProfile(ctx, gomock.Any()).
			Return(nil).
			Times(1)

		res, err := svc.UpdateProfile(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Updated Name", res.FirstName)
	})

	t.Run("profile not found", func(t *testing.T) {
		req := &profile.UpdateProfileRequest{
			AccountID: accountID,
		}

		mockRepo.EXPECT().
			GetProfileByAccountID(ctx, accountID).
			Return(nil, nil).
			Times(1)

		res, err := svc.UpdateProfile(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "profile not found")
	})
}
