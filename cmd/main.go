package cmd

import (
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/config"
	"github.com/FrostBitzX/smart-task-ai/internal/infrastructure/database"
)

func main() {
	cfg := config.NewConfig()
	db := database.NewDB(cfg)

	_ = db
	// [TODO]: Connect DB to repo
}
