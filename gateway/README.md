# Gateway Service

The gateway service acts as an API gateway and reverse proxy for all microservices in the home-server application. It routes incoming requests to the appropriate backend services using Docker Compose DNS.

## Architecture

```
gateway/
├── app/
│   └── main.go              # Entry point - initialization, route configuration and server startup
├── handlers/
│   └── handlers.go          # HTTP request handlers for each service route
├── services/
│   └── proxy.go             # Proxy logic and service discovery
├── static/
│   └── favicon.ico          # Static assets
├── config.yaml              # Service configuration
├── go.mod                   # Go module dependencies
└── Dockerfile               # Container build configuration
```

## Key Components

### main.go
- Initializes configuration and logging
- Sets up Gin router with middleware and routes
- Starts the HTTP server

### handlers/
Contains HTTP handlers for:
- Health checks
- Service-specific proxy handlers (profile, stats, camera, auth)
- Redirects and static content

### routes/
- **SetupRoutes()**: Registers all HTTP routes
- **SetupMiddleware()**: Configures CORS, recovery, and custom middleware
- Centralizes routing logic for easy maintenance

### services/
- **ProxyRequest()**: Generic proxy function that forwards requests to backend services
- **ServiceRegistry**: Maps service names to their ports
- Service discovery using Docker Compose DNS or environment variables

## Service Discovery

The gateway uses **Docker Compose DNS** for service discovery:

1. Service name in docker-compose.yml becomes the hostname
2. Example: `auth-service` is accessible at `http://auth-service:8080`
3. Environment variables can override defaults:
   - `AUTH_SERVICE_HOST` (default: service name)
   - `AUTH_SERVICE_PORT` (default: 8080)

## Routes

| Route Pattern | Target Service | Description |
|--------------|----------------|-------------|
| `/` | - | Redirects to `/profile` |
| `/health` | - | Gateway health check |
| `/profile/*` | profile | Profile service proxy |
| `/stats/*` | stats | Stats service proxy |
| `/camera/*` | camera | Camera service proxy |
| `/auth/*` | auth-service | Auth service proxy |

## Configuration

Edit `config.yaml` to configure:
- Service port
- Logging level and format
- Environment (development/production)

## Development

```bash
# Build the service
make build-gateway

# Start the gateway
make gateway-up

# View logs
make gateway-logs

# Restart after changes
make gateway-restart
```

## Features

- ✅ **Reverse Proxy**: Routes requests to backend microservices
- ✅ **Service Discovery**: Automatic service location via Docker DNS
- ✅ **CORS Support**: Configurable CORS middleware
- ✅ **Health Checks**: Built-in health endpoint
- ✅ **Structured Logging**: Using zap logger from common package
- ✅ **Error Handling**: Graceful error responses and recovery
- ✅ **Static Files**: Serves static assets (favicon, etc.)
- ✅ **Environment Aware**: Development and production modes

## Adding New Services

1. Add service to `services/proxy.go` ServiceRegistry
2. Create handler in `handlers/handlers.go`
3. Register route in `routes/routes.go`

Example:
```go
// In services/proxy.go
var ServiceRegistry = map[string]string{
    "new-service": "8080",
}

// In handlers/handlers.go
func NewServiceProxy(c *gin.Context) {
    services.ProxyRequest("new-service", c)
}

// In routes/routes.go (in SetupRoutes function)
router.Any("/newservice/*path", handlers.NewServiceProxy)
```

## Benefits of This Structure

- **Separation of Concerns**: Each package has a single responsibility
- **Maintainability**: Easy to find and modify specific functionality
- **Testability**: Handlers and services can be unit tested independently
- **Scalability**: Simple to add new routes and services
- **Readability**: Clean, organized code following Go best practices
- **Consistency**: Matches auth-service structure for familiarity
