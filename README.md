# Home Server

A microservices-based home server architecture built with Go, React, and PostgreSQL. This project provides a scalable foundation for managing various home automation and monitoring services.

## Architecture Overview

The home server consists of multiple microservices communicating through an API gateway:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   UI Service    │    │  Gateway Service │    │  Auth Service   │
│   (React)       │◄──►│   (Go/Gin)      │◄──►│   (Go)          │
│   Port: 3000    │    │   Port: 8080     │    │   Port: 8080    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                         │
                              ▼                         ▼
                    ┌─────────────────┐    ┌─────────────────┐
                    │  Consul Service │    │ PostgreSQL DB   │
                    │ (Service Disc.) │    │   Port: 5432    │
                    │   Port: 8500    │    │                 │
                    └─────────────────┘    └─────────────────┘
```

### Services

- **Gateway Service**: API gateway for routing requests to microservices
- **Auth Service**: Authentication and authorization service
- **UI Service**: React-based frontend application
- **PostgreSQL**: Primary database for all services
- **Consul**: Service discovery (optional, for Kubernetes deployment)

## Technology Stack

- **Backend**: Go 1.24.4 with Gin framework
- **Frontend**: React 18.3.1 with React Router
- **Database**: PostgreSQL 15
- **Service Discovery**: Consul (optional)
- **Containerization**: Docker & Docker Compose
- **Configuration**: YAML-based configuration files
- **Logging**: Structured logging with Zap

## Project Structure

```
home-server/
├── auth/                  # Authentication microservice
├── gateway/              # API gateway service
├── ui-service/           # React frontend application
├── discovery-service/    # Consul configuration
├── postgres/            # Database initialization scripts
├── common/              # Shared Go modules and utilities
├── docker-compose.yml   # Container orchestration
├── Makefile            # Build and deployment commands
└── go.work             # Go workspace configuration
```

## Prerequisites

- Docker and Docker Compose
- Python 3 (for database initialization scripts)
- Make (for build automation)
- Go 1.24+ (for local development)
- Node.js 18+ (for UI development)

## Quick Start

### 1. Environment Setup

Create `.env` files for each service with required passwords and configuration:

```bash
# Create environment files for services that need them
touch postgres/.env
touch auth/.env
touch gateway/.env
```

Add the following to `postgres/.env`:
```env
POSTGRES_USER=root_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=home_server
AUTH_DB_PASSWORD=your_auth_db_password
```

### 2. Start Services

```bash
# Generate database initialization scripts and start all services
make up

# Alternative: Run Docker Compose directly
docker compose up -d
```

### 3. Verify Installation

- **UI Service**: http://localhost:3000
- **Gateway API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Consul UI**: http://localhost:8500 (if enabled)

## Development

### Local Development Setup

1. **Set up Go workspace**:
   ```bash
   go work use ./auth ./gateway ./common
   ```

2. **Install dependencies**:
   ```bash
   # Go services
   cd auth && go mod tidy
   cd ../gateway && go mod tidy
   cd ../common && go mod tidy
   
   # React UI
   cd ../ui-service && npm install
   ```

3. **Run services locally**:
   ```bash
   # Start database first
   docker compose up postgres -d
   
   # Run individual services
   cd auth && go run app/main.go
   cd gateway && go run app/main.go
   cd ui-service && npm start
   ```

### Database Management

```bash
# Connect to PostgreSQL
make postgres-login

# Clean generated files
make clean
```

### Available Make Commands

- `make up` - Start all services with Docker Compose
- `make down` - Stop all services
- `make generate-init` - Generate database initialization scripts
- `make postgres-login` - Connect to PostgreSQL container
- `make clean` - Remove generated files

## Configuration

Each service uses YAML configuration files:

- **Auth Service**: `auth/config.yaml`
- **Gateway**: `gateway/config.yaml`
- **Common**: `common/config/config_template.yaml`

Configuration includes:
- Service ports and names
- Database connections
- Logging levels and formats
- API timeouts and retry policies
- Security settings

## API Documentation

### Gateway Endpoints

The gateway service proxies requests to microservices:

- `GET /profile/*` - Profile service endpoints
- `GET /stats/*` - Statistics service endpoints  
- `GET /camera/*` - Camera service endpoints
- `GET /health` - Health check endpoint

### Authentication

The auth service provides user authentication and authorization using JWT tokens.

## Monitoring and Health Checks

- Health check endpoints available at `/health` for each service
- Structured logging with configurable levels
- Service discovery with Consul for production deployments

## Production Deployment

### Docker Compose (Recommended for single-node)
```bash
make up
```

### Kubernetes (For multi-node clusters)
- Remove Consul dependency (K8s provides service discovery)
- Use Kubernetes Services and Ingress
- Configure persistent volumes for PostgreSQL

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **Database connection errors**: Ensure PostgreSQL is running and `.env` files are configured
2. **Service discovery issues**: Check Consul configuration and network connectivity
3. **Port conflicts**: Verify no other services are using the required ports

### Logs

View logs for specific services:
```bash
docker compose logs -f [service-name]
docker compose logs -f postgres
docker compose logs -f gateway
```
