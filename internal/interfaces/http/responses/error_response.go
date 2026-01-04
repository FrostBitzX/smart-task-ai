package responses

import (
	"github.com/FrostBitzX/smart-task-ai/internal/errors/apperrors"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents the standard error response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   ErrorDetail `json:"error"`
}

// SuccessResponse represents the standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(c *fiber.Ctx, err error) error {
	var status int
	var code string
	var message string

	if appErr, ok := apperrors.IsAppError(err); ok {
		status = appErr.Status
		code = appErr.Code
		message = appErr.Message
	} else {
		status = apperrors.StatusCode(err)
		code = "INTERNAL_SERVER_ERROR"
		message = err.Error()
	}

	errorResponse := ErrorResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Error: ErrorDetail{
			Code:    status,
			Message: code,
		},
	}

	return c.Status(status).JSON(errorResponse)
}

func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	successResponse := SuccessResponse{
		Success: true,
		Message: msg,
		Data:    data,
		Error:   nil,
	}

	return c.Status(fiber.StatusOK).JSON(successResponse)
}
