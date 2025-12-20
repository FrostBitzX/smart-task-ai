package profile

type CreateProfileRequest struct {
	FirstName        string `json:"first_name" validate:"min=3,max=20"`
	LastName           string `json:"last_name" validate:"min=3,max=20"`
	Nickname           string `json:"nickname" validate:"min=3,max=20"`
	AvatarPath           string `json:"avatar_path"`
}

type CreateProfileResponse struct {
	AccountID string `json:"account_id"`
	ProfileID string `json:"profile_id"`
}
