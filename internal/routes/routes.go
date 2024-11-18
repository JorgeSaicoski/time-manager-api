package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
	"github.com/JorgeSaicoski/time-manager-api/internal/database"
	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
)

func SetupRouter(cfg *database.Config) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{
					"Origin",
					"Content-Type",
					"Accept",
					"Authorization",
					"X-Requested-With",
			},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:          12 * time.Hour,  // Preflight cache duration
	}))
	userHandler := handlers.NewUserHandler(cfg.DB)
	totalTimeHandler := handlers.NewTotalTimerHandler(cfg.DB)

	SetupUserRoutes(router, userHandler)
	SetupTotalTimeRoutes(router, totalTimeHandler)

	return router
}
