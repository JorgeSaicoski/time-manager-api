package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/database"
	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
)

func RegisterUserRoutes(r *gin.Engine, cfg *database.Config) {
	userHandler := handlers.NewUserHandler(cfg.DB)
	users := r.Group("/users")
	{
		users.POST("", userHandler.CreateUser)       // POST   /users
		users.GET("", userHandler.GetUsers)          // GET    /users
		users.GET("/:id", userHandler.GetUser)       // GET    /users/:id
		users.PUT("/:id", userHandler.UpdateUser)    // PUT    /users/:id
		users.DELETE("/:id", userHandler.DeleteUser) // DELETE /users/:id
	}
}
