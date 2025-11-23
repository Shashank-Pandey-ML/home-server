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

		// Stats service routes (proxied to stats-service)
		api.Any("/stats", handlers.StatsServiceProxy)

		// Auth service routes (all under /auth/*)
		api.Any("/auth/*path", handlers.AuthServiceProxy)

		// File service routes (protected - future)
		api.Any("/files/*path", handlers.FileServiceProxy)

		// Camera service routes (protected)
		api.Any("/camera/*path", handlers.CameraServiceProxy)
	}

	// Serve static files (CSS, JS, images) from React build
	router.Static("/static", "./ui-build/static")
	router.StaticFile("/favicon.ico", "./ui-build/favicon.ico")
	router.StaticFile("/manifest.json", "./ui-build/manifest.json")
	router.StaticFile("/logo192.png", "./ui-build/logo192.png")
	router.StaticFile("/logo512.png", "./ui-build/logo512.png")
	router.StaticFile("/robots.txt", "./ui-build/robots.txt")

	// Serve React app for all non-API routes (client-side routing)
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
