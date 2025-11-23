package main

import (
	"fmt"

	"stats/handlers"

	"github.com/gin-gonic/gin"
	"github.com/shashank/home-server/common/config"
	"github.com/shashank/home-server/common/logging"
	"go.uber.org/zap"
)

// init initializes the stats service configuration and logger
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
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API routes - All backend microservices under /api/v1
	api := router.Group(config.AppConfig.API.BaseURL)
	{
		// Stats endpoint
		api.GET("/stats", handlers.StatsHandler)
	}

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
