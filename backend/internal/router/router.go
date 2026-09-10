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

	// Services
	authService := service.NewAuthService(userRepo, cfg)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)

	api := r.Group("/api")
	{
		// Public routes
		api.POST("/auth/login", authHandler.Login)

		// Authenticated routes
		authed := api.Group("/")
		authed.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			authed.POST("/auth/change-password", authHandler.ChangePassword)
		}
	}

	return r
}