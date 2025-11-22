# Authentication Service

A Go-based microservice providing authentication and authorization functionality for the home server ecosystem. Built with clean architecture principles and comprehensive logging.

## Overview

The Auth Service is responsible for:
- User authentication and authorization
- JWT token management
- User profile management
- Database interactions for user data
- Health monitoring and logging

## Technology Stack

- **Language**: Go 1.24.4
- **Framework**: Standard Go libraries with custom HTTP handlers
- **Database**: PostgreSQL with dedicated `auth` database
- **Logging**: Structured logging with Uber Zap
- **Configuration**: YAML-based configuration with Viper
- **Containerization**: Docker support

## Project Structure

```
auth-service/
├── app/
│   ├── main.go           # Application entry point
├── config.yaml           # Service configuration
├── Dockerfile            # Container configuration
├── go.mod               # Go module dependencies
├── go.sum               # Dependency checksums
└── README.md            # This file
```

## Configuration

The service uses `config.yaml` for configuration management:

### Database Setup

The database is automatically initialized with:
- Database: `auth`
- User: `auth_user`
- Permissions: Full access to `auth` database

## Environment Variables

Create a `.env` file in the auth-service directory:

```env
# Database configuration
AUTH_DB_PASSWORD=your_secure_password

# Optional: Override config values
AUTH_SERVICE_PORT=8080
AUTH_LOG_LEVEL=info
AUTH_ENVIRONMENT=prod
```

## API Endpoints

### Health Check
- **GET** `/health` - Service health status

### Authentication (Planned)
- **POST** `/api/v1/auth/login` - User login
- **POST** `/api/v1/auth/logout` - User logout
- **POST** `/api/v1/auth/refresh` - Token refresh

### User Management (Planned)
- **POST** `/api/v1/users` - Create user
- **GET** `/api/v1/users/{id}` - Get user by ID
- **PUT** `/api/v1/users/{id}` - Update user
- **DELETE** `/api/v1/users/{id}` - Delete user

## Development Setup

### Prerequisites
- Go 1.24+
- PostgreSQL 15+
- Docker (optional)

### Local Development

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Set up environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start PostgreSQL** (if not using Docker):
   ```bash
   # Start from project root
   docker compose up postgres -d
   ```

4. **Run the service**:
   ```bash
   go run app/main.go
   ```

5. **Build binary**:
   ```bash
   go build -o app/main app/main.go
   ```

### Docker Development

1. **Build container**:
   ```bash
   docker build -t auth-service .
   ```

2. **Run with Docker Compose** (from project root):
   ```bash
   docker compose up auth-service -d
   ```

## Testing

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
go test -tags=integration ./...
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Logging

The service uses structured logging with configurable levels:

- **Debug**: Detailed debugging information
- **Info**: General operational messages
- **Warn**: Warning conditions
- **Error**: Error conditions

Log format can be configured as `json` or `text` in the configuration file.

### Log Examples

```json
{
  "level": "info",
  "ts": "2025-07-02T10:00:00.000Z",
  "caller": "main.go:25",
  "msg": "Auth service started successfully",
  "service": "auth",
  "port": 8080
}
```

## Security Considerations

- Passwords are never included in JSON responses
- TLS encryption support (configurable)
- Database credentials via environment variables
- JWT token-based authentication
- Role-based access control

## Dependencies

### Direct Dependencies
- `github.com/shashank/home-server/common` - Shared utilities and models
- `go.uber.org/zap` - Structured logging

### Indirect Dependencies
- Viper for configuration management
- PostgreSQL driver
- Various utility libraries

## Monitoring

### Health Endpoint
The service provides a health check endpoint that reports:
- Service status
- Database connectivity
- Configuration validation

### Metrics (Planned)
- Request/response times
- Error rates
- Active user sessions
- Database connection pool status

## Deployment

### Docker Deployment
```bash
# Build and run
docker build -t auth-service .
docker run -d \
  --name auth-service \
  -p 8080:8080 \
  --env-file .env \
  auth-service
```

### Production Considerations
- Use environment-specific configuration files
- Enable TLS in production
- Set up proper database backup procedures
- Configure log rotation and monitoring
- Use secrets management for sensitive data

## Contributing

1. Follow Go best practices and conventions
2. Add tests for new functionality
3. Update documentation for API changes
4. Use structured logging for all operations
5. Maintain backward compatibility when possible

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   ```bash
   # Check database is running
   docker compose ps postgres
   
   # Verify credentials
   docker compose exec postgres psql -U auth_user -d auth
   ```

2. **Configuration Issues**
   ```bash
   # Validate YAML syntax
   go run app/main.go --validate-config
   ```

3. **Port Conflicts**
   ```bash
   # Check if port is in use
   lsof -i :8080
   ```

### Debug Mode
Enable debug logging by setting:
```yaml
logging:
  level: "debug"
```

## Future Enhancements

- [ ] Complete JWT authentication implementation
- [ ] Add password hashing and validation
- [ ] Implement rate limiting
- [ ] Add OAuth2 provider support
- [ ] Implement audit logging
- [ ] Add user session management
- [ ] Create admin dashboard endpoints
- [ ] Add password reset functionality
