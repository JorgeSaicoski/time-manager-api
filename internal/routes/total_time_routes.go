package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/JorgeSaicoski/time-manager-api/internal/handlers"
    "github.com/JorgeSaicoski/time-manager-api/internal/middleware"
)

func SetupTotalTimeRoutes(router *gin.Engine, totalTimeHandler *handlers.TotalTimeHandler) {
    protected := router.Group("/totaltime")
    protected.Use(middleware.AuthMiddleware())
    protected.Use(middleware.VerifyUserMiddleware())

    protected.GET("/health", func(c *gin.Context) {
        userID := c.GetString("user_id")
        if userID == "" {
            c.JSON(401, gin.H{
                "status":  "error",
                "message": "Unauthorized access",
            })
            return
        }
        c.JSON(200, gin.H{
            "status":  "ok",
            "message": "Health check passed",
        })
    })

    protected.POST("/", totalTimeHandler.CreateTotalTime)
    protected.PUT("/", totalTimeHandler.CloseTotalTime)
    protected.GET("/", totalTimeHandler.GetTotalTime)
}
