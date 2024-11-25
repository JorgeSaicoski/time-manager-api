package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/JorgeSaicoski/time-manager-api/internal/middleware"
	"github.com/JorgeSaicoski/time-manager-api/internal/models"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" binding:"email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UpdateUserRoleRequest struct {
	IsSystemAdmin bool `json:"is_system_admin"`
}

type UserResponse struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	IsSystemAdmin bool      `json:"is_system_admin"`
	CreatedAt     time.Time `json:"created_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User

	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		IsSystemAdmin: user.IsSystemAdmin,
		CreatedAt:     user.CreatedAt,
	})
}

func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Email != "" {
		var existingUser models.User
		if err := h.db.Where("email = ? AND id != ?", req.Email, userID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already taken"})
			return
		}
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		IsSystemAdmin: user.IsSystemAdmin,
		CreatedAt:     user.CreatedAt,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := h.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if err := h.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Model(&models.User{}).Where("id = ?", userID).Update("is_system_admin", req.IsSystemAdmin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if result := h.db.Where("email = ?", req.Email).First(&existingUser); !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if result := h.db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, user.IsSystemAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Name:          user.Name,
			IsSystemAdmin: user.IsSystemAdmin,
			CreatedAt:     user.CreatedAt,
		},
		"tokens": tokens,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := h.db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokens, err := middleware.GenerateTokenPair(user.ID, user.Email, user.IsSystemAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Name:          user.Name,
			IsSystemAdmin: user.IsSystemAdmin,
			CreatedAt:     user.CreatedAt,
		},
		"tokens": tokens,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := &middleware.JWTClaims{}
	token, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return middleware.JwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	tokens, err := middleware.GenerateTokenPair(claims.UserID, claims.Email, claims.IsSystemAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	requestingUserID, _ := c.Get("user_id")
	isAdmin, _ := c.Get("is_system_admin")

	if requestingUserID != userID && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	tx := h.db.Begin()

	if err := tx.Where("user_id = ?", userID).Delete(&models.TotalTime{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.WorkTime{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := tx.Where("owner_id = ?", userID).Delete(&models.Project{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := tx.Delete(&models.User{}, userID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []models.User

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var total int64
	h.db.Model(&models.User{}).Count(&total)

	if err := h.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Name:          user.Name,
			IsSystemAdmin: user.IsSystemAdmin,
			CreatedAt:     user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userResponses,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  (int(total) + limit - 1) / limit,
			"total_items":  total,
		},
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		IsSystemAdmin: user.IsSystemAdmin,
		CreatedAt:     user.CreatedAt,
	})
}

func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.deleteUser(c, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *UserHandler) deleteUser(_ *gin.Context, userID int64) error {
	tx := h.db.Begin()

	if err := tx.Where("user_id = ?", userID).Delete(&models.TotalTime{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.WorkTime{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("owner_id = ?", userID).Delete(&models.Project{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(&models.User{}, userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
