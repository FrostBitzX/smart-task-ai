package requests

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ParseAndValidate[T any](c *fiber.Ctx) (*T, error) {
	var body T

	if err := c.BodyParser(&body); err != nil {
		return nil, c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Operation failed",
			"data":    nil,
			"error": fiber.Map{
				"code":    400,
				"message": "Invalid JSON format",
			},
		})
	}

	if err := validate.Struct(&body); err != nil {
		return nil, validationErrorResponse(c, err)
	}

	return &body, nil
}

func ParseAndValidateQuery[T any](c *fiber.Ctx) (*T, error) {
	var q T

	if err := c.QueryParser(&q); err != nil {
		return nil, c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Operation failed",
			"data":    nil,
			"error": fiber.Map{
				"code":    400,
				"message": "Invalid query parameters",
			},
		})
	}

	if err := validate.Struct(&q); err != nil {
		return nil, validationErrorResponse(c, err)
	}

	return &q, nil
}

func ParseAndValidateDetailed[T any](c *fiber.Ctx) (*T, map[string]string, error) {
	var body T
	if err := c.BodyParser(&body); err != nil {
		return nil, nil, err
	}

	if err := validate.Struct(body); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			details := make(map[string]string)
			for _, e := range errs {
				switch e.Tag() {
				case "required":
					details[e.Field()] = "is required"
				case "email":
					details[e.Field()] = "must be a valid email"
				case "min":
					details[e.Field()] = "must be at least " + e.Param() + " characters"
				case "max":
					details[e.Field()] = "must be at most " + e.Param() + " characters"
				default:
					details[e.Field()] = "is invalid"
				}
			}
			return nil, details, err
		}
		return nil, nil, err
	}

	return &body, nil, nil
}

func validationErrorResponse(c *fiber.Ctx, err error) error {
	msg := "Validation failed"

	if errs, ok := err.(validator.ValidationErrors); ok {
		fe := errs[0]
		switch fe.Tag() {
		case "required":
			msg = fe.Field() + " is required"
		case "email":
			msg = "email format is invalid"
		case "min":
			msg = fe.Field() + " must be at least " + fe.Param() + " characters"
		case "max":
			msg = fe.Field() + " must be at most " + fe.Param() + " characters"
		default:
			msg = fe.Field() + " is invalid"
		}
	}

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Operation failed",
		"data":    nil,
		"error": fiber.Map{
			"code":    400,
			"message": msg,
		},
	})
}
