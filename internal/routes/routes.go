package routes

import (
	"time"

	"github.com/JorgeSaicoski/time-manager-api/internal/database"
	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
	"github.com/JorgeSaicoski/time-manager-api/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *database.Config) *gin.Engine {
	if cfg == nil || cfg.DB == nil {
		panic("database configuration cannot be nil")
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	repo := repository.New(cfg.DB)
	userHandler := handlers.NewUserHandler(cfg.DB)
	totalTimeHandler := handlers.NewTotalTimerHandler(repo)

	SetupUserRoutes(router, userHandler)
	SetupTotalTimeRoutes(router, totalTimeHandler)
	return router
}
