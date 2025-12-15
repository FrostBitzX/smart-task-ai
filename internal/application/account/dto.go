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

type AccountDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

type ListAccountsRequest struct {
	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type PaginationDTO struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

type ListAccountsResponse struct {
	Items      []AccountDTO  `json:"items"`
	Pagination PaginationDTO `json:"pagination"`
}
