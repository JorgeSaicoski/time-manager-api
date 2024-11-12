package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/JorgeSaicoski/time-manager-api/internal/database"
)

func SetupRouter(cfg *database.Config) *gin.Engine {
	router := gin.Default()

	RegisterUserRoutes(router, cfg)

	return router
}
