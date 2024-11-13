package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
)

func SetupUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	router.POST("/auth/register", userHandler.Register)
	router.POST("/auth/login", userHandler.Login)
	router.POST("/auth/refresh", userHandler.RefreshToken)

	// Protected routes
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/me", userHandler.GetCurrentUser)
		protected.PUT("/users/me", userHandler.UpdateCurrentUser)
		protected.POST("/users/me/change-password", userHandler.ChangePassword)
		protected.DELETE("/users/me", userHandler.DeleteCurrentUser)

		// Admin routes
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
