package middlewares

import (
	"strings"
	"time"

	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v2"
)

func FiberLoggerMiddleware(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Path() == "/healthz" || strings.HasPrefix(c.Path(), "/api/v2/healthz") {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		log.Info("http_request", map[string]interface{}{
			"method":  c.Method(),
			"path":    c.Path(),
			"status":  c.Response().StatusCode(),
			"latency": duration.String(),
			"ip":      c.IP(),
		})

		if err != nil {
			log.Error("http_error", map[string]interface{}{
				"error": err.Error(),
				"path":  c.Path(),
			})
		}

		return err
	}
}
