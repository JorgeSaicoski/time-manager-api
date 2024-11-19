package routes

import (
	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupTotalTimeRoutes(router *gin.Engine, totalTimeHandler *handlers.TotalTimeHandler) {
	// Base group with auth middleware
	protected := router.Group("/totaltime")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.VerifyUserMiddleware())

	// User-specific routes with user_id parameter
	protected.GET("/user/:user_id/health", func(c *gin.Context) {
		userID := c.Param("user_id")
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Health check passed",
			"user_id": userID,
		})
	})

	protected.POST("/user/:user_id", totalTimeHandler.CreateTotalTime)
	protected.PUT("/user/:user_id", totalTimeHandler.CloseTotalTime)
	protected.GET("/user/:user_id", totalTimeHandler.GetTotalTime)
	//get :totaltime_id/user/:user_id/
	//get /user/:user_id/query (will take by month/week/day)
	//get :totaltime_id/user/:user_id/
	//put :totaltime_id/user/:user_id/

}
