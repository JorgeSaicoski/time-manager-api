package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS_TOKEN_DURATION  = time.Hour * 24      // 24 hours
	REFRESH_TOKEN_DURATION = time.Hour * 24 * 30 // 30 days
)

var JwtSecret []byte

func init() {
	secret := os.Getenv("SECRET-JWT")
	if secret == "" {
		panic("JWT secret is not set in environment variables")
	}
	JwtSecret = []byte(secret)
}

type JWTClaims struct {
	UserID        int64  `json:"user_id"`
	Email         string `json:"email"`
	IsSystemAdmin bool   `json:"is_system_admin"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (j *JWTClaims) Valid() error {
	panic("unimplemented")
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateTokenPair(userID int64, email string, isSystemAdmin bool) (*TokenPair, error) {
	accessClaims := JWTClaims{
		UserID:        userID,
		Email:         email,
		IsSystemAdmin: isSystemAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ACCESS_TOKEN_DURATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "time-manager-api",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(JwtSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := JWTClaims{
		UserID:        userID,
		Email:         email,
		IsSystemAdmin: isSystemAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(REFRESH_TOKEN_DURATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "time-manager-api",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(JwtSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenStr := bearerToken[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return JwtSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("is_system_admin", claims.IsSystemAdmin)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isSystemAdmin, exists := c.Get("is_system_admin")
		if !exists || !isSystemAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func VerifyUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user_id from JWT (set by AuthMiddleware)
		jwtUserIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
			c.Abort()
			return
		}

		jwtUserID, ok := jwtUserIDInterface.(int64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			c.Abort()
			return
		}

		// Extract user_id from URL parameter
		paramUserIDStr := c.Param("user_id")
		if paramUserIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
			c.Abort()
			return
		}

		// Convert user_id from string to int64
		paramUserID, err := strconv.ParseInt(paramUserIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
			c.Abort()
			return
		}

		// Compare the two user IDs
		if paramUserID != jwtUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
