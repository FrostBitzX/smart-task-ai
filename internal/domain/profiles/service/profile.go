package service

import (
	"context"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/profile"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
)

type ProfileService struct {
	repo profiles.ProfileRepository
}

func NewProfileService(repo profiles.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) CreateProfile(ctx context.Context, req *profile.CreateProfileRequest) (*entity.Profile, error) {
	if req == nil {
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// create domain entity
	now := time.Now()
	prof := &entity.Profile{
		ID:         uuid.New(),
		AccountID:  req.AccountID,
		FirstName:  lo.ToPtr(req.FirstName),
		LastName:   lo.ToPtr(req.LastName),
		Nickname:   lo.ToPtr(req.Nickname),
		AvatarPath: lo.ToPtr(req.AvatarPath),
		State:      "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// persist account to database
	if err := s.repo.CreateProfile(ctx, prof); err != nil {
		return nil, apperrors.NewInternalServerError("failed to create profile", "CREATE_PROFILE_ERROR", err)
	}

	return prof, nil
}
