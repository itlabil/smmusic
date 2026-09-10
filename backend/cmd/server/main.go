package main

import (
	"log"

	"github.com/itlabil/smmusic/backend/internal/config"
	"github.com/itlabil/smmusic/backend/internal/router"
)

func main() {
	cfg := config.Load()
	db := config.NewDatabase(cfg)
	defer db.Close()

	r := router.Setup(db, cfg)

	log.Printf("Server starting on port %s (env: %s)", cfg.Port, cfg.Env)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}