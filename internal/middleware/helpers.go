package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserRequesting(c *gin.Context) (int64, error) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("user ID not found in context or invalid format")
	}
	return userID, nil
}
