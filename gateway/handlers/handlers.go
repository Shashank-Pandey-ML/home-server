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

// ServeReactApp returns middleware that serves the React SPA
// In production: serves static files from build directory
// In development: can proxy to React dev server on port 3000
func ServeReactApp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Option 1: Proxy to React dev server (development mode)
		// Uncomment this block when running React with `npm start`
		/*
			services.ProxyRequest("ui-service", c)
			return
		*/

		// Option 2: Serve React build files (production mode)
		// Serve static files from ./ui-build directory
		// This will be the React production build output
		path := c.Request.URL.Path

		// Check if file exists in ui-build directory
		// If not, serve index.html to support client-side routing
		if _, err := http.Dir("./ui-build").Open(path); err != nil {
			// File not found, serve index.html for React Router
			c.File("./ui-build/index.html")
			return
		}

		// File exists, serve it
		c.FileFromFS(path, http.Dir("./ui-build"))
	}
}
