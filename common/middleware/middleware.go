package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/shashank/home-server/common/config"
	"github.com/shashank/home-server/common/logging"
)

// CorsMiddleware handles Cross-Origin Resource Sharing (CORS) headers
// Uses allowed_origins from JWT config
func CorsMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get allowed origins from JWT config
		allowedOrigins := config.AppConfig.JWT.AllowedOrigins

		// If no origins configured, allow all (not recommended for production)
		if len(allowedOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// Check if the request origin is in the allowed list
			origin := c.GetHeader("Origin")
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// RateLimitMiddleware implements rate limiting to prevent abuse
// Default: 100 requests per minute, can be configured
func RateLimitMiddleware() gin.HandlerFunc {
	// Create a rate limiter: 100 requests per minute
	// TODO: Make this configurable via config file
	limiter := rate.NewLimiter(rate.Every(time.Minute/100), 100)

	return gin.HandlerFunc(func(c *gin.Context) {
		// In production, you might want per-IP rate limiting using a map of limiters
		if !limiter.Allow() {
			logging.Log.Warn("Rate limit exceeded",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("service", config.AppConfig.Service.Name))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// SecurityHeadersMiddleware adds security-related HTTP headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Prevent information leakage - use service name from config
		serviceName := config.AppConfig.Service.Name
		c.Header("Server", serviceName+"-service")

		// Strict transport security (only for HTTPS)
		if c.Request.TLS != nil || config.AppConfig.Security.EnableTLS {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// Content security policy
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	})
}

// RequestLoggingMiddleware logs all incoming requests with structured logging
func RequestLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			logging.Log.Info("HTTP Request",
				zap.String("method", param.Method),
				zap.String("path", param.Path),
				zap.Int("status", param.StatusCode),
				zap.Duration("latency", param.Latency),
				zap.String("ip", param.ClientIP),
				zap.String("user_agent", param.Request.UserAgent()),
				zap.String("service", config.AppConfig.Service.Name),
			)
			return ""
		},
		Output: nil, // We handle logging ourselves
	})
}

// HealthCheckMiddleware provides a simple health check response
// This can be used by services that don't need custom health logic
func HealthCheckMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.URL.Path == config.AppConfig.Health.Endpoint {
			c.JSON(http.StatusOK, gin.H{
				"status":    "healthy",
				"service":   config.AppConfig.Service.Name,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"version":   "1.0.0",
			})
			c.Abort()
			return
		}
		c.Next()
	})
}
