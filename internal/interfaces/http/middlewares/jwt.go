package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		// ตรวจสอบว่า Bearer มี prefix
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(auth, bearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
		}

		tokenStr := strings.TrimPrefix(auth, bearerPrefix)

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// Extract all claims
		username, _ := claims["Username"].(string)
		accountID, _ := claims["AccountId"].(string)
		email, _ := claims["Email"].(string)

		if accountID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing AccountId in token",
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
