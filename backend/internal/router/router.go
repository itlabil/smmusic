package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/itlabil/smmusic/backend/internal/config"
	"github.com/itlabil/smmusic/backend/internal/handler"
	"github.com/itlabil/smmusic/backend/internal/middleware"
	"github.com/itlabil/smmusic/backend/internal/repository"
	"github.com/itlabil/smmusic/backend/internal/service"
)

func Setup(db *sql.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	songRepo := repository.NewSongRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	transcodeService := service.NewTranscodeService(songRepo, 2) // max 2 concurrent transcode workers
	songService := service.NewSongService(songRepo, cfg, transcodeService)
	interactionService := service.NewInteractionService(interactionRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	songHandler := handler.NewSongHandler(songService)
	interactionHandler := handler.NewInteractionHandler(interactionService)

	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/login", authHandler.Login)

		// Authenticated routes
		authed := api.Group("/")
		authed.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			authed.POST("/auth/change-password", authHandler.ChangePassword)

			authed.POST("/songs/upload", songHandler.Upload)
			authed.GET("/songs", songHandler.List)
			authed.GET("/songs/:id/stream", songHandler.Stream)

			authed.POST("/songs/:id/like", interactionHandler.Like)
			authed.DELETE("/songs/:id/like", interactionHandler.Unlike)
			authed.GET("/songs/liked", interactionHandler.ListLiked)

			authed.POST("/songs/:id/play", interactionHandler.RecordPlay)
			authed.GET("/songs/recently-played", interactionHandler.ListRecentlyPlayed)

			// Admin-only routes
			admin := authed.Group("/admin")
			admin.Use(middleware.AdminRequired())
			{
				admin.POST("/users", authHandler.CreateUser)
				admin.GET("/users", authHandler.ListUsers)
			}
		}
	}

	return r
}