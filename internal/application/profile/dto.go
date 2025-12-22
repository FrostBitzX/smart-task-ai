package profile

import "time"

type CreateProfileRequest struct {
	AccountID  string  `json:"account_id"`
	FirstName  string  `json:"first_name" validate:"min=0,max=20"`
	LastName   string  `json:"last_name" validate:"min=0,max=20"`
	Nickname   string  `json:"nickname" validate:"min=0,max=20"`
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
	FirstName  string    `json:"first_name" validate:"min=0,max=20"`
	LastName   string    `json:"last_name" validate:"min=0,max=20"`
	Nickname   string    `json:"nickname" validate:"min=0,max=20"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	State      string    `json:"state"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
