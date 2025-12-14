package main

import (
	"log"
	"os"

	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/logger"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/middlewares"
	"github.com/FrostBitzX/smart-task-ai/internal/interfaces/http/routes"
	"gorm.io/gorm"

	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/config"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Application struct {
	Server *fiber.App
	Config *config.Config
	DB     *gorm.DB
	Logger logger.Logger
}

func main() {
	app := fiber.New()
	cfg := config.NewConfig()
	db := database.NewDB(cfg)

	app := fiber.New(fiber.Config{
		AppName: "smart-task-ai",
	})

	zapLogger := logger.NewZapLogger()
	app.Use(middlewares.FiberLoggerMiddleware(zapLogger))

	routes.RegisterPublicRoutes(app, db, zapLogger)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("🚀 server running on %s", addr)
	log.Fatal(app.Listen(addr))
}
