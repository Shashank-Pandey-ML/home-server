package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gateway/handlers"
	gateway_middleware "gateway/middleware"

	"github.com/shashank/home-server/common/config"
	"github.com/shashank/home-server/common/logging"
	"github.com/shashank/home-server/common/middleware"
)

// init initializes the gateway service configuration and logger
func init() {
	// Load the configuration
	if err := config.LoadConfig("config.yaml"); err != nil {
		// Log the error and panic if configuration loading fails
		// This ensures that the application does not start with an invalid configuration.
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize the logger with the loaded configuration
	if err := logging.InitLogger(config.AppConfig.Logging, config.AppConfig.Service.Name); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	logging.Log.Info("Gateway service initialization completed successfully")
}

func main() {
	// Set Gin mode based on environment
	if config.AppConfig.Service.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.Default()

	// Add middleware
	router.Use(middleware.RequestLoggingMiddleware())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	// Configure trusted proxies
	router.SetTrustedProxies(nil)

	// Health check endpoint (no /api prefix for gateway health)
	router.GET("/health", handlers.HealthHandler)

	// API routes - All backend microservices under /api/v1
	api := router.Group(config.AppConfig.API.BaseURL)
	{
		// Conditional auth middleware - skips auth for specific public endpoints
		api.Use(gateway_middleware.ConditionalAuthMiddleware([]string{
			"/api/v1/auth/login",
			// "/api/v1/auth/public-key",
		}))

		// Auth service routes (all under /auth/*)
		api.Any("/auth/*path", handlers.AuthServiceProxy)

		// Stats service routes (protected)
		api.Any("/stats/*path", handlers.StatsServiceProxy)

		// File service routes (protected - future)
		api.Any("/files/*path", handlers.FileServiceProxy)

		// Camera service routes (protected)
		api.Any("/camera/*path", handlers.CameraServiceProxy)
	}

	// Serve React UI - must be last to catch all non-API routes
	// Profile page is public (default landing page)
	router.GET("/", handlers.ServeReactApp())
	router.GET("/profile", handlers.ServeReactApp())
	router.GET("/profile/*path", handlers.ServeReactApp())

	// All other UI routes require authentication
	router.Use(gateway_middleware.OptionalAuthMiddleware())
	router.NoRoute(handlers.ServeReactApp())

	// Start the server
	port := fmt.Sprintf(":%d", config.AppConfig.Service.Port)
	logging.Log.Info("Starting gateway service",
		zap.String("port", port),
		zap.String("environment", config.AppConfig.Service.Environment),
	)

	if err := router.Run(port); err != nil {
		logging.Log.Fatal("Failed to start gateway service", zap.Error(err))
	}
}
