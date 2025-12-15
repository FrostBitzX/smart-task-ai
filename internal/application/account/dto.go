package account

type CreateAccountRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Email           string `json:"email"`
}

type CreateAccountResponse struct {
	AccountID string `json:"account_id"`
}
