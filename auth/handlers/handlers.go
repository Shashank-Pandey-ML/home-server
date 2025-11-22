package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shashank/home-server/auth-service/services"
	"github.com/shashank/home-server/common/db"
	"github.com/shashank/home-server/common/logging"
	"github.com/shashank/home-server/common/models"
)

// LoginRequest represents the JSON payload for login requests
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginResponse represents the JSON response for successful login
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// UserResponse represents user data in API responses (without sensitive info)
type UserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

// RefreshRequest represents the JSON payload for token refresh requests
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest represents the JSON payload for logout requests
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// loginHandler handles user authentication and JWT token generation
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Log.Warn("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Log the login attempt (without password)
	logging.Log.Info("Login attempt", zap.String("email", req.Email))

	accessToken, refreshToken, expiresIn, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		logging.Log.Warn("Login failed", zap.String("email", req.Email), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Log successful login
	logging.Log.Info("User logged in successfully", zap.String("email", req.Email))

	// Return tokens and user info
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	})
}

// logoutHandler handles user logout and token invalidation
func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	var req LogoutRequest

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Log.Warn("Invalid logout request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Extract user info from JWT context (set by middleware)
	userIdStr, exists := c.Get("user_id")
	if !exists {
		logging.Log.Warn("Logout attempt without valid user context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication context",
		})
		return
	}

	userID, _ := strconv.ParseUint(userIdStr.(string), 10, 64)
	err := h.authService.Logout(c.Request.Context(), req.RefreshToken, uint(userID))
	if err != nil {
		logging.Log.Error("Failed to logout user",
			zap.String("user_id", userIdStr.(string)),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout user",
		})
		return
	}

	logging.Log.Info("User logged out successfully", zap.String("user_id", userIdStr.(string)))

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// refreshHandler handles JWT token refresh using refresh tokens
func (h *AuthHandler) RefreshHandler(c *gin.Context) {
	var req RefreshRequest

	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.Log.Warn("Invalid refresh request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// TODO: Validate refresh token against database
	user, err := h.authService.ValidateRefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		logging.Log.Warn("Invalid refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	// Generate new access token (optionally new refresh token too)
	accessToken, refreshToken, expiresIn, err := services.GenerateTokenPair(user)
	if err != nil {
		logging.Log.Error("Failed to generate new tokens", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh authentication tokens",
		})
		return
	}

	// TODO: Update refresh token in database (optional: rotate refresh tokens)

	logging.Log.Info("Tokens refreshed successfully", zap.Uint("user_id", user.ID))

	// Return new tokens
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	})
}

// getUserProfileHandler returns the current user's profile information
func (h *AuthHandler) GetUserProfileHandler(c *gin.Context) {
	// Extract user info from JWT context (set by middleware)
	userIdStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication context",
		})
		return
	}

	userID, _ := strconv.ParseUint(userIdStr.(string), 10, 64)
	user, err := h.authService.GetUserByID(c.Request.Context(), uint(userID))
	if err != nil {
		logging.Log.Error("Failed to fetch user profile",
			zap.String("user_id", userIdStr.(string)),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch user profile",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:      userIdStr.(string),
		Email:   user.Email,
		Name:    user.Name,
		IsAdmin: user.IsAdmin,
	})
}

// updateUserProfileHandler updates the current user's profile information
func (h *AuthHandler) UpdateUserProfileHandler(c *gin.Context) {
	// Extract user info from JWT context (set by middleware)
	userIdStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication context",
		})
		return
	}

	var userUpdate models.User
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		logging.Log.Warn("Invalid profile update request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	userID, _ := strconv.ParseUint(userIdStr.(string), 10, 64)
	userUpdate.ID = uint(userID)

	if err := h.authService.UpdateUserProfile(c.Request.Context(), &userUpdate); err != nil {
		logging.Log.Error("Failed to update user profile",
			zap.String("user_id", userIdStr.(string)),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user profile",
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:      userIdStr.(string),
		Email:   userUpdate.Email,
		Name:    userUpdate.Name,
		IsAdmin: userUpdate.IsAdmin,
	})
}

// getPublicKeyHandler provides the JWT public key for token validation
func (h *AuthHandler) GetPublicKeyHandler(c *gin.Context) {
	publicKeyPEM, err := h.authService.GetPublicKeyPEM(c.Request.Context())
	if err != nil {
		logging.Log.Error("Failed to get public key", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve public key",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"public_key": publicKeyPEM,
		"algorithm":  "RS256",
		"key_type":   "RSA",
	})
}

// HealthCheckHandler checks the health of the auth service
type HealthCheckHandler struct {
	db *db.DB
}

// NewHealthCheckHandler creates a new HealthCheckHandler
func NewHealthCheckHandler(db *db.DB) *HealthCheckHandler {
	return &HealthCheckHandler{
		db: db,
	}
}

// healthCheckHandler provides a health check endpoint
func (h *HealthCheckHandler) HealthCheckHandler(c *gin.Context) {
	// TODO: Add database connectivity check
	status := gin.H{
		"status":    "healthy",
		"service":   "auth",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	databaseHealth := h.db.HealthCheck(c.Request.Context())

	status["database"] = databaseHealth
	c.JSON(http.StatusOK, status)
}
