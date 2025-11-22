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
		// Auth service routes (public - no auth required)
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.AuthServiceProxy)
			auth.GET("/public-key", handlers.AuthServiceProxy)
		}

		// Protected routes - require authentication
		protected := api.Group("")
		protected.Use(gateway_middleware.AuthMiddleware())
		{
			// User management routes (protected)
			users := protected.Group("/users")
			{
				users.POST("/logout", handlers.AuthServiceProxy)
				users.POST("/refresh", handlers.AuthServiceProxy)
				users.GET("/profile", handlers.AuthServiceProxy)
				users.PUT("/profile", handlers.AuthServiceProxy)
			}

			// Stats service routes (protected)
			protected.Any("/stats/*path", handlers.StatsServiceProxy)

			// File service routes (protected - future)
			protected.Any("/files/*path", handlers.FileServiceProxy)

			// Camera service routes (protected)
			protected.Any("/camera/*path", handlers.CameraServiceProxy)
		}
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
