package responses

import (
	appError "github.com/FrostBitzX/smart-task-ai/pkg/apperror"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents the standard error response
type ErrorResponse struct {
	Code    string `json:"code" example:"example code"`
	Message string `json:"message" example:"example message"`
}

func Error(c *fiber.Ctx, err error) error {
	var code string
	if appErr, ok := err.(*appError.AppError); ok {
		code = appErr.Code
	} else {
		code = "UNKNOWN"
	}

	status := fiber.StatusInternalServerError
	if appErr, ok := err.(*appError.AppError); ok && appErr.HTTPStatus != 0 {
		status = appErr.HTTPStatus
	} else {
		status = appError.StatusCode(err)
	}

	return c.Status(status).JSON(ErrorResponse{
		Code:    code,
		Message: err.Error(),
	})
}
