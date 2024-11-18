package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
)

func SetupTotalTimeRoutes(router *gin.Engine, totalTimeHandler *handlers.TotalTimeHandler) {
  protected := router.Group("/totaltime")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.VerifyUserMiddleware())
	{
		//health
		protected.GET("/health", func(c *gin.Context) {
			fmt.Print("Params:")
			fmt.Println(c.Params)
			fmt.Print("keys:")
			fmt.Println(c.Keys)
			fmt.Print("request:")
			fmt.Println(c.Request)
				c.JSON(200, gin.H{
						"status": "ok",
						"message": "Health check passed",
				})
		})
		// Totaltime routes
    protected.POST("/create", totalTimeHandler.CreateTotalTime)
    protected.PUT("/close", totalTimeHandler.CloseTotalTime)
	}
}
