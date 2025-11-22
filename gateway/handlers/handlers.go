package handlers

import (
	"net/http"
	"time"

	"gateway/services"

	"github.com/gin-gonic/gin"
)

// HealthHandler returns the health status of the gateway
func HealthHandler(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"service":   "gateway",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0.0",
	}

	c.JSON(http.StatusOK, status)
}

// RedirectToProfile redirects the root path to /profile
func RedirectToProfile(c *gin.Context) {
	c.Redirect(http.StatusFound, "/profile")
}

// ProfileServiceProxy proxies requests to the profile service
func ProfileServiceProxy(c *gin.Context) {
	services.ProxyRequest("profile", c)
}

// StatsServiceProxy proxies requests to the stats service
func StatsServiceProxy(c *gin.Context) {
	services.ProxyRequest("stats", c)
}

// CameraServiceProxy proxies requests to the camera service
func CameraServiceProxy(c *gin.Context) {
	services.ProxyRequest("camera", c)
}

// AuthServiceProxy proxies requests to the auth service
func AuthServiceProxy(c *gin.Context) {
	services.ProxyRequest("auth", c)
}

// FileServiceProxy proxies requests to the file service (future)
func FileServiceProxy(c *gin.Context) {
	services.ProxyRequest("file-service", c)
}

// ServeReactApp serves the React SPA from the build directory
func ServeReactApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Serve index.html for all routes (client-side routing)
		c.File("./ui-build/index.html")
	}
}
