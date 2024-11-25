package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	router.GET("/health", func(c *gin.Context) {
		fmt.Println("test")
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Health check passed",
		})
	})

	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)
	router.POST("/auth/refresh", userHandler.RefreshToken)

	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users/me", userHandler.GetCurrentUser)
		protected.PUT("/users/me", userHandler.UpdateCurrentUser)
		protected.POST("/users/me/change-password", userHandler.ChangePassword)
		protected.DELETE("/users/me", userHandler.DeleteCurrentUser)

		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.GET("/users/:id", userHandler.GetUser)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
			admin.POST("/users/:id/role", userHandler.UpdateUserRole)
		}
	}
}
