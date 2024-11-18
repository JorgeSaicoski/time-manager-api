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
		protected.GET("/health/:user_id", func(c *gin.Context) {
        userID, exists := c.Params.Get("user_id")
        if !exists {
					c.JSON(404, gin.H{
							"status":  "error",
							"message": "Health check passed with user",
					})
					return
        }

        fmt.Printf("user id: %s\n", userID)

				c.JSON(200, gin.H{
            "status":  "ok",
            "message": "Health check passed with user",
        })
    })
		// Totaltime routes
    protected.POST("/create", totalTimeHandler.CreateTotalTime)
    protected.PUT("/close", totalTimeHandler.CloseTotalTime)
	}
}
