package middlewares

import (
	"os"
	"strings"

	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Success: false,
				Message: "missing token",
				Data:    nil,
				Error: responses.ErrorDetail{
					Code:    fiber.StatusUnauthorized,
					Message: "MISSING_TOKEN",
				},
			})
		}

		// ตรวจสอบว่า Bearer มี prefix
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(auth, bearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Success: false,
				Message: "invalid token format",
				Data:    nil,
				Error: responses.ErrorDetail{
					Code:    fiber.StatusUnauthorized,
					Message: "INVALID_TOKEN_FORMAT",
				},
			})
		}

		tokenStr := strings.TrimPrefix(auth, bearerPrefix)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Success: false,
				Message: "invalid token",
				Data:    nil,
				Error: responses.ErrorDetail{
					Code:    fiber.StatusUnauthorized,
					Message: "INVALID_TOKEN",
				},
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Success: false,
				Message: "invalid token claims",
				Data:    nil,
				Error: responses.ErrorDetail{
					Code:    fiber.StatusUnauthorized,
					Message: "INVALID_TOKEN_CLAIMS",
				},
			})
		}

		// Extract all claims
		username, _ := claims["Username"].(string)
		accountID, _ := claims["AccountId"].(string)
		email, _ := claims["Email"].(string)

		if accountID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(responses.ErrorResponse{
				Success: false,
				Message: "missing account id in token",
				Data:    nil,
				Error: responses.ErrorDetail{
					Code:    fiber.StatusUnauthorized,
					Message: "MISSING_ACCOUNT_ID",
				},
			})
		}

		// Save all claims into context (Locals)
		c.Locals("jwt_claims", map[string]interface{}{
			"AccountId": accountID,
			"Email":     email,
			"Username":  username,
		})

		return c.Next()
	}
}
