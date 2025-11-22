package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shashank/home-server/common/logging"
	"github.com/shashank/home-server/common/models"
	"go.uber.org/zap"
)

// Cache for the RSA public key
var (
	cachedPublicKey     *rsa.PublicKey
	publicKeyCacheMutex sync.RWMutex
	publicKeyExpiry     time.Time
)

// AuthMiddleware validates JWT tokens by calling the auth
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logging.Log.Debug("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Check for Bearer token format
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logging.Log.Debug("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token locally using public key
		claims, err := validateJWTLocally(token)
		if err != nil {
			logging.Log.Warn("Token validation failed",
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Store claims in context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("is_admin", claims.IsAdmin)

		logging.Log.Debug("Token validated successfully",
			zap.String("user_id", claims.UserID),
			zap.String("email", claims.Email),
		)

		c.Next()
	}
}

// validateJWTLocally validates JWT token using cached public key from auth-service
func validateJWTLocally(tokenString string) (*models.JWTClaims, error) {
	// Get public key (from cache or fetch)
	publicKey, err := getPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Validate token type (should be "access" for API requests)
	if claims.Type != models.TokenTypeAccess {
		return nil, errors.New("invalid token type, expected access token")
	}

	return claims, nil
}

// getPublicKey retrieves the RSA public key from cache or fetches it from auth-service
func getPublicKey() (*rsa.PublicKey, error) {
	// Check cache first (with read lock)
	publicKeyCacheMutex.RLock()
	if cachedPublicKey != nil && time.Now().Before(publicKeyExpiry) {
		publicKeyCacheMutex.RUnlock()
		return cachedPublicKey, nil
	}
	publicKeyCacheMutex.RUnlock()

	// Cache miss or expired, fetch new key (with write lock)
	publicKeyCacheMutex.Lock()
	defer publicKeyCacheMutex.Unlock()

	// Double-check after acquiring write lock
	if cachedPublicKey != nil && time.Now().Before(publicKeyExpiry) {
		return cachedPublicKey, nil
	}

	// Fetch public key from auth-service
	authServiceURL := getAuthServiceURL() + "/api/v1/auth/public-key"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(authServiceURL)
	if err != nil {
		logging.Log.Error("Failed to fetch public key from auth-service",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch public key: %s", string(body))
	}

	// Parse response
	var keyResp models.PublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&keyResp); err != nil {
		return nil, fmt.Errorf("failed to parse public key response: %w", err)
	}

	// Parse PEM-encoded public key
	block, _ := pem.Decode([]byte(keyResp.PublicKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	// Parse RSA public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	// Cache the public key for 1 hour
	cachedPublicKey = rsaPublicKey
	publicKeyExpiry = time.Now().Add(1 * time.Hour)

	logging.Log.Info("Public key fetched and cached",
		zap.Time("expires_at", publicKeyExpiry),
	)

	return cachedPublicKey, nil
}

// getAuthServiceURL returns the auth-service URL using Docker Compose DNS
func getAuthServiceURL() string {
	// In Docker Compose, service name is the hostname
	return "http://auth-service:8080"
}

// OptionalAuthMiddleware validates JWT tokens if present, but doesn't require them
// Useful for routes that behave differently based on authentication status
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without setting user context
			c.Next()
			return
		}

		// Try to validate token
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			token := tokenParts[1]
			claims, err := validateJWTLocally(token)
			if err == nil {
				// Valid token, set user context
				c.Set("user_id", claims.UserID)
				c.Set("email", claims.Email)
				c.Set("is_admin", claims.IsAdmin)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}
