package profile

import "time"

type CreateProfileRequest struct {
	AccountID  string  `json:"account_id"`
	FirstName  string  `json:"first_name" validate:"required"`
	LastName   string  `json:"last_name" validate:"required"`
	Nickname   *string `json:"nickname,omitempty"`
	AvatarPath *string `json:"avatar_path,omitempty"`
}

type CreateProfileResponse struct {
	AccountID string `json:"account_id"`
	ProfileID string `json:"profile_id"`
}

type GetProfileByAccountIDRequest struct {
	AccountID string `json:"account_id"`
}

type GetProfileByAccountIDResponse struct {
	AccountID  string    `json:"account_id"`
	FirstName  string    `json:"first_name" validate:"min=4,max=50"`
	LastName   string    `json:"last_name" validate:"min=4,max=50"`
	Nickname   *string   `json:"nickname,omitempty" validate:"omitempty,max=20"`
	AvatarPath *string   `json:"avatar_path,omitempty" validate:"omitempty"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	AccountID  string  `json:"account_id"`
	FirstName  string  `json:"first_name" validate:"required,min=4,max=50"`
	LastName   string  `json:"last_name" validate:"required,min=4,max=50"`
	Nickname   *string `json:"nickname,omitempty" validate:"omitempty,max=20"`
	AvatarPath *string `json:"avatar_path,omitempty" validate:"omitempty"`
}

type UpdateProfileResponse struct {
	AccountID  string  `json:"account_id"`
	ProfileID  string  `json:"profile_id"`
	FirstName  string  `json:"first_name" validate:"min=4,max=50"`
	LastName   string  `json:"last_name" validate:"min=4,max=50"`
	Nickname   *string `json:"nickname,omitempty" validate:"omitempty,max=20"`
	AvatarPath *string `json:"avatar_path,omitempty" validate:"omitempty"`
}
