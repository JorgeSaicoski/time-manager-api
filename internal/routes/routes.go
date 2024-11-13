package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/database"
	"github.com/JorgeSaicoski/time-manager-api/internal/handlers"
)

func SetupRouter(cfg *database.Config) *gin.Engine {
	router := gin.Default()

	userHandler := handlers.NewUserHandler(cfg.DB)

	SetupUserRoutes(router, userHandler)

	return router
}
