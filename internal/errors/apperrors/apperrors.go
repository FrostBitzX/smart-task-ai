package apperrors

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/logger"
)

// AppError represents a domain/application-level error that can be mapped to HTTP responses.
type AppError struct {
	Status   int         // HTTP status code
	Code     string      // Machine-readable error code (e.g. "INTERNAL_SERVER_ERROR")
	Message  string      // Human-readable message
	Details  interface{} // Optional additional details for response
	RawError error       // Underlying error for logging
}

func (e *AppError) Error() string {
	if e.RawError != nil {
		return e.Message + ": " + e.RawError.Error()
	}
	return e.Message
}

// ------------------------
// Factory methods
// ------------------------

func NewBadRequestError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewUnauthorizedError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusUnauthorized,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewForbiddenError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusForbidden,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewNotFoundError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusNotFound,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewConflictError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusConflict,
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewInternalServerError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// ------------------------
// Common Errors (Migrated from pkg/apperror)
// ------------------------

var (
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrInvalidData    = errors.New("invalid data")
	ErrRecordNotFound = logger.ErrRecordNotFound
)

// StatusCode maps generic errors to HTTP status codes
func StatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status
	}

	switch {
	case errors.Is(err, ErrUnauthorized):
		return fiber.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return fiber.StatusForbidden
	case errors.Is(err, ErrRecordNotFound):
		return fiber.StatusNotFound
	case errors.Is(err, ErrInvalidData):
		return fiber.StatusBadRequest
	case errors.Is(err, ErrInternalServer):
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}

// IsAppError tries to cast a generic error into *AppError.
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
