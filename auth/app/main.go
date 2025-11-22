package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/shashank/home-server/auth/handlers"
	auth_middleware "github.com/shashank/home-server/auth/middleware"
	"github.com/shashank/home-server/auth/services"
	"github.com/shashank/home-server/common/config"
	"github.com/shashank/home-server/common/db"
	"github.com/shashank/home-server/common/logging"
	"github.com/shashank/home-server/common/middleware"
	"github.com/shashank/home-server/common/models"
)

// init initializes the gateway service configuration and logger
func init() {
	// Load the configuration
	if err := config.LoadConfig("config.yaml"); err != nil {
		// Log the error and panic if configuration loading fails
		// This ensures that the application does not start with an invalid configuration.
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize the logger with the configuration loaded from config.yaml
	if err := logging.InitLogger(config.AppConfig.Logging, config.AppConfig.Service.Name); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	logging.Log.Info("Auth service initialization completed successfully")

	// Initialize JWT keys after configuration is loaded
	if err := services.InitializeJWTKeys(); err != nil {
		logging.Log.Fatal("Failed to initialize JWT keys", zap.Error(err))
	}
}

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Initialize database connection
	database, err := db.InitDbConnection(config.AppConfig.Database, logging.Log)
	if err != nil {
		logging.Log.Fatal("Failed to initialize database connection", zap.Error(err))
	}
	defer database.Close()

	database.AutoMigrate(&models.User{}) // Ensure User model is migrated

	// Build dependencies
	healthCheckHandler := handlers.NewHealthCheckHandler(database)
	userRepo := db.NewUserRepository(database)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Add middleware
	router.Use(middleware.RequestLoggingMiddleware())
	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.RateLimitMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())

	// Health check endpoint
	router.GET("/health", healthCheckHandler.HealthCheckHandler)

	// Authentication routes
	api := router.Group(config.AppConfig.API.BaseURL)
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.LoginHandler)

			// Public key endpoint for gateway service
			auth.GET("/public-key", authHandler.GetPublicKeyHandler)
		}

		// User management routes (protected)
		users := api.Group("/users")
		users.Use(auth_middleware.JwtAuthMiddleware())
		{
			auth.POST("/logout", authHandler.LogoutHandler)
			auth.POST("/refresh", authHandler.RefreshHandler)
			users.GET("/profile", authHandler.GetUserProfileHandler)
			users.PUT("/profile", authHandler.UpdateUserProfileHandler)
		}
	}

	// Start the server
	port := fmt.Sprintf(":%d", config.AppConfig.Service.Port)
	logging.Log.Info("Starting auth service", zap.String("port", port),
		zap.String("environment", config.AppConfig.Service.Environment))

	if err := router.Run(port); err != nil {
		logging.Log.Fatal("Failed to start auth service", zap.Error(err))
	}
}
