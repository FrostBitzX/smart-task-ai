package models

type CreateAccount struct {
	Username        string
	Password        string
	ConfirmPassword string
	Email           string
}
