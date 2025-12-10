package main

import (
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/config"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	cfg := config.NewConfig()
	db := database.NewDB(cfg)

	_ = db
	// [TODO]: Connect DB to repo

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running...")
	})

	auth := app.Group("/auth")
	_ = auth

	app.Listen(":3000")
}
