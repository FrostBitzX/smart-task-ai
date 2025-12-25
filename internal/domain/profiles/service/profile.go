package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/profile"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
)

type ProfileService struct {
	repo profiles.ProfileRepository
}

func NewProfileService(repo profiles.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

func (s *ProfileService) GetProfileByAccountID(ctx context.Context, accountID string) (*entity.Profile, error) {
	prof, err := s.repo.GetProfileByAccountID(ctx, accountID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to get profile by account id", "GET_PROFILE_BY_ACCOUNT_ID_ERROR", err)
	}

	return prof, nil
}

func (s *ProfileService) CreateProfile(ctx context.Context, req *profile.CreateProfileRequest) (*entity.Profile, error) {
	if req == nil {
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// Check if profile already created
	exists, err := s.GetProfileByAccountID(ctx, req.AccountID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to get profile by account id", "GET_PROFILE_BY_ACCOUNT_ID_ERROR", err)
	}
	if exists != nil {
		return nil, apperrors.NewBadRequestError("profile already exists", "PROFILE_ALREADY_EXISTS", nil)
	}

	// create domain entity
	now := time.Now()
	prof := &entity.Profile{
		ID:         uuid.New(),
		AccountID:  uuid.MustParse(req.AccountID),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Nickname:   req.Nickname,
		AvatarPath: req.AvatarPath,
		State:      "active",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// persist account to database
	err = s.repo.CreateProfile(ctx, prof)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperrors.NewBadRequestError(
				fmt.Sprintf("profile with account ID %s already exists", req.AccountID),
				"PROFILE_ALREADY_EXISTS",
				err,
			)
		}
		return nil, apperrors.NewInternalServerError("failed to create profile", "CREATE_PROFILE_ERROR", err)
	}

	return prof, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, req *profile.UpdateProfileRequest) (*entity.Profile, error) {
	if req == nil {
		return nil, apperrors.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// Check if profile already created
	exists, err := s.GetProfileByAccountID(ctx, req.AccountID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("failed to get profile by account id", "GET_PROFILE_BY_ACCOUNT_ID_ERROR", err)
	}
	if exists == nil {
		return nil, apperrors.NewBadRequestError("profile not found", "PROFILE_NOT_FOUND", nil)
	}

	// create domain entity
	now := time.Now()
	prof := &entity.Profile{
		ID:         exists.ID,
		AccountID:  uuid.MustParse(req.AccountID),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Nickname:   req.Nickname,
		AvatarPath: req.AvatarPath,
		State:      "active",
		CreatedAt:  exists.CreatedAt,
		UpdatedAt:  now,
	}

	// persist account to database
	err = s.repo.UpdateProfile(ctx, prof)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperrors.NewBadRequestError(
				fmt.Sprintf("profile with account ID %s already exists", req.AccountID),
				"PROFILE_ALREADY_EXISTS",
				err,
			)
		}
		return nil, apperrors.NewInternalServerError("failed to update profile", "UPDATE_PROFILE_ERROR", err)
	}

	return prof, nil
}
