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
	"github.com/joho/godotenv"
)

type Application struct {
	Server *fiber.App
	Config *config.Config
	DB     *gorm.DB
	Logger logger.Logger
}

func main() {
	_ = godotenv.Load()

	cfg := config.NewConfig()
	db := database.NewDB(cfg)

	app := fiber.New(fiber.Config{
		AppName: "smart-task-ai",
	})

	zapLogger := logger.NewZapLogger()
	app.Use(middlewares.FiberLoggerMiddleware(zapLogger))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,https://smart-task-ai-fe-prod.vercel.app",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	routes.RegisterPublicRoutes(app, db, zapLogger)
	routes.RegisterPrivateRoutes(app, db, zapLogger)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("🚀 server running on %s", addr)
	log.Fatal(app.Listen(addr))
}
