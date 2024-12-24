package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gateway/app/services"
)

// RegisterRoutes sets up all the routes for the gateway
func RegisterRoutes(router *gin.Engine) {
	// Middleware
	// router.Use(AuthMiddleware())

	// Proxy in this context refes to a reverse proxy. If your application is behind a reverse
	// proxy (e.g., Nginx, Cloudflare Tunnel, AWS ALB), the proxy will forward the client request
	// to your backend service.
	// TODO: Using Cloudflare Tunnel: Your backend will receive traffic forwarded by Cloudflare. Configure
	// Cloudflare's IPs as trusted proxies.
	router.SetTrustedProxies(nil)

	// Redirect the root route to the /profile page
	router.GET("/", func(c *gin.Context) {
		// Redirect to the profile page
		c.Redirect(http.StatusFound, "/profile")
	})

	// Registering the
	router.StaticFile("/favicon.ico", "./app/static/favicon.ico")

	// Registeting the health API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Gateway is healthy",
		})
	})

	// Route groups
	router.GET("/profile/*path", ProfileServiceProxy)
	router.GET("/stats/*path", StatsServiceProxy)
	router.GET("/camera/*path", CameraServiceProxy)
}

// Function to proxy request to profile service
func ProfileServiceProxy(c *gin.Context) {
	proxyRequest("profile", c)
}

// Function to proxy request to profile service
func StatsServiceProxy(c *gin.Context) {
	proxyRequest("stats", c)
}

// Function to proxy request to profile service
func CameraServiceProxy(c *gin.Context) {
	proxyRequest("camera", c)
}

// Generic Function to proxy request to respective service
func proxyRequest(serviceName string, c *gin.Context) {
	path := c.Param("path")
	method := c.Request.Method
	response, err := services.ProxyRequest(serviceName, path, method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Copy response back to client
	c.DataFromReader(response.StatusCode, response.ContentLength, response.Header.Get("Content-Type"), response.Body, nil)
}
