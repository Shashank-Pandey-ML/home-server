package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"

	"gateway/app/config"
	"gateway/app/logger"
	"gateway/app/routes"
)

func init() {
	logger.Logger.Info("Starting gateway service")
}

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Register routes
	routes.RegisterRoutes(router)

	// Start the server
	port := strconv.Itoa(config.AppConfig.Service.Port)
	logger.Logger.Sugar().Infof("Gateway is running on port: %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}
