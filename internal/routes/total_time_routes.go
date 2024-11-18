package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
)

func SetupTotalTimeRoutes(router *gin.Engine, totalTimeHandler *handlers.TotalTimeHandler) {
  protected := router.Group("/totaltime")
	protected.Use(middleware.AuthMiddleware())
	{
		//health
		protected.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
		// Totaltime routes
    protected.POST("/create", totalTimeHandler.CreateTotalTime)
    protected.PUT("/close", totalTimeHandler.CloseTotalTime)
	}
}
