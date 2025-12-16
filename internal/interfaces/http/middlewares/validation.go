package middlewares

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CreateAccountRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=20"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=4"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=4,eqfield=Password"`
}

var validate = validator.New()

func ValidateCreateAccountRequest[T any](c *fiber.Ctx) error {
	var req T

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Operation failed",
			"data":    nil,
			"error": fiber.Map{
				"code":    400,
				"message": "Invalid JSON format",
			},
		})
	}

	if err := validate.Struct(&req); err != nil {
		return handleValidationError(c, err)
	}

	// Store validated request in context for handler to use
	c.Locals("validatedRequest", req)
	return c.Next()
}

func handleValidationError(c *fiber.Ctx, err error) error {
	var errorMessage string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessage = getValidationErrorMessage(validationErrors[0])
	} else {
		errorMessage = "Validation failed"
	}

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Operation failed",
		"data":    nil,
		"error": fiber.Map{
			"code":    400,
			"message": errorMessage,
		},
	})
}

func getValidationErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()

	switch tag {
	case "required":
		return strings.ToLower(field) + " is required"
	case "min":
		return strings.ToLower(field) + " must be at least " + fe.Param() + " characters"
	case "max":
		return strings.ToLower(field) + " must be at most " + fe.Param() + " characters"
	case "email":
		return "email format is invalid"
	case "eqfield":
		return "password and confirm_password must match"
	default:
		return strings.ToLower(field) + " is invalid"
	}
}
