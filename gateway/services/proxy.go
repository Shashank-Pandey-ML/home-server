package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shashank/home-server/common/logging"
	"go.uber.org/zap"
)

// ServiceRegistry holds the service name to port mapping
var ServiceRegistry = map[string]string{
	"profile":      "8080",
	"stats":        "8080",
	"camera":       "8080",
	"auth-service": "8080",
	"file-service": "8080",
	"ui-service":   "3000", // React dev server or production build server
}

// ProxyRequest forwards the incoming request to the target service
func ProxyRequest(serviceName string, c *gin.Context) {
	// Get target service configuration
	targetHost := getServiceHost(serviceName)
	targetPort := getServicePort(serviceName)

	// Build the target URL using Docker Compose DNS
	path := c.Param("path")
	if path == "" {
		path = "/"
	}
	targetURL := fmt.Sprintf("http://%s:%s%s", targetHost, targetPort, path)

	logging.Log.Debug("Proxying request",
		zap.String("service", serviceName),
		zap.String("method", c.Request.Method),
		zap.String("target_url", targetURL),
	)

	// Create a new HTTP request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		logging.Log.Error("Failed to create proxy request",
			zap.Error(err),
			zap.String("service", serviceName),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create proxy request",
		})
		return
	}

	// Copy headers from original request
	copyHeaders(req.Header, c.Request.Header)

	// Forward the request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		logging.Log.Error("Proxy request failed",
			zap.Error(err),
			zap.String("service", serviceName),
			zap.String("target_url", targetURL),
		)
		c.JSON(http.StatusBadGateway, gin.H{
			"error": fmt.Sprintf("Service %s is unavailable", serviceName),
		})
		return
	}
	defer resp.Body.Close()

	// Copy response headers back to client
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set response status and stream body
	c.Status(resp.StatusCode)
	if _, err := io.Copy(c.Writer, resp.Body); err != nil {
		logging.Log.Error("Failed to stream response body",
			zap.Error(err),
			zap.String("service", serviceName),
		)
	}
}

// getServiceHost returns the hostname for the service
// Checks environment variable first, then falls back to service name (Docker Compose DNS)
func getServiceHost(serviceName string) string {
	envKey := strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")) + "_SERVICE_HOST"
	if host := os.Getenv(envKey); host != "" {
		return host
	}
	// Docker Compose DNS: service name is the hostname
	return serviceName
}

// getServicePort returns the port for the service
// Checks environment variable first, then registry, then defaults to 8080
func getServicePort(serviceName string) string {
	envKey := strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")) + "_SERVICE_PORT"
	if port := os.Getenv(envKey); port != "" {
		return port
	}
	if port, exists := ServiceRegistry[serviceName]; exists {
		return port
	}
	return "8080" // default port
}

// copyHeaders copies HTTP headers from source to destination
func copyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
