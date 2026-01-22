package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/application/profile"
	"github.com/FrostBitzX/smart-task-ai/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles"
	"github.com/FrostBitzX/smart-task-ai/internal/domain/profiles/entity"
)

type ProfileService struct {
	repo        profiles.ProfileRepository
	accountRepo accounts.AccountRepository
}

func NewProfileService(repo profiles.ProfileRepository, accountRepo accounts.AccountRepository) *ProfileService {
	return &ProfileService{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (s *ProfileService) CheckAndGetProfile(ctx context.Context, accountID string) (*entity.Profile, error) {
	prof, err := s.repo.CheckAndGetProfile(ctx, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperror.NewInternalServerError("failed to check and get profile", "CHECK_AND_GET_PROFILE_ERROR", nil)
	}

	return prof, nil
}

func (s *ProfileService) CreateProfile(ctx context.Context, req *profile.CreateProfileRequest) (*entity.Profile, error) {
	if req == nil {
		return nil, apperror.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// Check if account exists
	_, err := s.accountRepo.GetAccount(ctx, req.AccountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewBadRequestError(
				"account ID does not exist",
				"ACCOUNT_NOT_FOUND",
				nil,
			)
		}
		return nil, apperror.NewInternalServerError("failed to check account existence", "CHECK_ACCOUNT_ERROR", nil)
	}

	// Check if profile already exists
	exists, err := s.CheckAndGetProfile(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}
	if exists != nil {
		return nil, apperror.NewBadRequestError("profile already exists", "PROFILE_ALREADY_EXISTS", nil)
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

	// persist profile to database
	err = s.repo.CreateProfile(ctx, prof)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to create profile", "CREATE_PROFILE_ERROR", nil)
	}

	return prof, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, req *profile.UpdateProfileRequest) (*entity.Profile, error) {
	if req == nil {
		return nil, apperror.NewBadRequestError("invalid request body", "INVALID_REQUEST", nil)
	}

	// Check if account exists
	_, err := s.accountRepo.GetAccount(ctx, req.AccountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewBadRequestError(
				fmt.Sprintf("account with ID %s does not exist", req.AccountID),
				"ACCOUNT_NOT_FOUND",
				nil,
			)
		}
		return nil, apperror.NewInternalServerError("failed to check account existence", "CHECK_ACCOUNT_ERROR", nil)
	}

	// Check if profile exists
	exists, err := s.CheckAndGetProfile(ctx, req.AccountID)
	if err != nil {
		return nil, err
	}
	if exists == nil {
		return nil, apperror.NewBadRequestError("profile not found", "PROFILE_NOT_FOUND", nil)
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
		return nil, apperror.NewInternalServerError("failed to update profile", "UPDATE_PROFILE_ERROR", nil)
	}

	return prof, nil
}
