package requests

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ParseAndValidate[T any](c *fiber.Ctx) (*T, error) {
	var body T

	if err := c.BodyParser(&body); err != nil {
		return nil, err
	}

	if err := validate.Struct(body); err != nil {
		return nil, err
	}

	return &body, nil
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
