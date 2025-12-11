package http

import (
	"context"
	"net/http"

	"github.com/FrostBitzX/smart-task-ai/internal/domain/accounts"
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"
)

func (h *AccountHandler) CreateAccount(ctx context.Context, req CreateAccountRequestObject) (CreateAccountResponseObject, error) {
	// Prepare data
	acc := AccountRequest{
		Username:        req.Body.Username,
		Email:           string(req.Body.Email),
		Password:        req.Body.Password,
		ConfirmPassword: req.Body.ConfirmPassword,
	}

	// Call the service layer
	accResponse, err := h.AccountService.CreateAccount(ctx, acc)
	if err != nil {
		appErr, ok := apperrors.IsAppError(err)
		if ok {
			return handleCreateAccountAppError(appErr), nil
		}

		return CreateAccount500JSONResponse{
			NewInternalServerError("Failed to create account", "INTERNAL_SERVER_ERROR", err.Error()),
		}, nil
	}

	// Create API response
	account := accounts.FromAccountModel(&accResponse.Account)
	apiAccountResponse := toAPIAccountResponse(*account)

	// Return successful response
	return CreateAccount201JSONResponse{
		Data:    apiAccountResponse,
		Message: "Create account successfully",
		Success: true,
		Error:   nil,
	}, nil
}

func handleCreateAccountAppError(err *apperrors.AppError) CreateAccountResponseObject {
	switch err.Status {
	case http.StatusBadRequest:
		return CreateAccount400JSONResponse{
			NewBadRequestError(err.Message, err.Code, err.Details),
		}
	case http.StatusConflict:
		return CreateAccount409JSONResponse{
			NewConflictError(err.Message, err.Code, err.Details),
		}
	case http.StatusInternalServerError:
		return CreateAccount500JSONResponse{
			NewInternalServerError(err.Message, err.Code, err.Details),
		}
	}
	return CreateAccount500JSONResponse{
		NewInternalServerError(err.Message, err.Code, err.Details),
	}
}
