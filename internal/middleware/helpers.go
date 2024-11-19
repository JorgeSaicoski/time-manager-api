package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetUserRequesting(c *gin.Context) (int64, error) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type")
	}

	return userID, nil
}
