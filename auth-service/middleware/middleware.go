package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shashank/home-server/auth-service/services"
	"github.com/shashank/home-server/common/logging"
)

// jwtAuthMiddleware validates JWT tokens and extracts user information
func JwtAuthMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logging.Log.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check for Bearer token format
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logging.Log.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := services.ValidateJWTToken(token)
		if err != nil {
			logging.Log.Warn("Invalid JWT token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract user information from claims and set in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_is_admin", claims.IsAdmin)

		c.Next()
	})
}
