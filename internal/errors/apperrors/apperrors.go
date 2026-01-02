package apperrors

import "net/http"

// AppError represents a domain/application-level error that can be mapped to HTTP responses.
type AppError struct {
	Status  int         // HTTP status code
	Code    string      // Machine-readable error code (e.g. "INTERNAL_SERVER_ERROR")
	Message string      // Human-readable message
	Details interface{} // Optional additional details (can be string, map, struct, etc.)
}

func (e *AppError) Error() string {
	return e.Message
}

// NewBadRequestError creates an AppError representing a 400 Bad Request.
func NewBadRequestError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewConflictError creates an AppError representing a 409 Conflict.
func NewConflictError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusConflict,
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

// NewInternalServerError creates an AppError representing a 500 Internal Server Error.
func NewInternalServerError(message, code string, details interface{}) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// IsAppError tries to cast a generic error into *AppError.
// It returns the casted *AppError and a boolean indicating whether the cast succeeded.
func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}

	appErr, ok := err.(*AppError)
	if ok {
		return appErr, true
	}

	return nil, false
}
