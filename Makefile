# Paths
POSTGRES_INIT_DIR=postgres
PYTHON_SCRIPT=generate_postgres_setup.py

# Go services
GO_SERVICES=auth gateway stats common

# Default target
.PHONY: all
all: build-ui up!

# Help target - show available commands
.PHONY: help
help:
	@echo "📚 Home Server Makefile Commands"
	@echo ""
	@echo "🚀 Quick Start:"
	@echo "  make all              - Build UI and start all services (default)"
	@echo "  make up               - Start all services (without rebuild)"
	@echo "  make up!              - Build UI and start all services (force rebuild)"
	@echo "  make down             - Stop all services (keeps data)"
	@echo ""
	@echo "🔨 Build Commands:"
	@echo "  make build            - Build all services"
	@echo "  make build-ui         - Build React UI"
	@echo "  make build-gateway    - Build gateway service"
	@echo "  make build-auth       - Build auth service"
	@echo "  make build-stats      - Build stats service"
	@echo "  make build-info       - Show which services can be built"
	@echo ""
	@echo "🐳 Service Management:"
	@echo "  make <service>-up     - Start specific service (e.g., make postgres-up)"
	@echo "  make <service>-down   - Stop specific service"
	@echo "  make <service>-restart - Restart specific service"
	@echo "  make <service>-logs   - View logs for specific service"
	@echo ""
	@echo "📦 Go Module Management:"
	@echo "  make tidy             - Run go mod tidy on all Go services"
	@echo "  make deps             - Download Go dependencies"
	@echo "  make verify           - Verify Go modules"
	@echo ""
	@echo "👤 User Management:"
	@echo "  make create-user      - Create new user (interactive)"
	@echo "  make list-users       - List all users"
	@echo ""
	@echo "🧹 Cleanup Commands:"
	@echo "  make clean            - Remove generated files and binaries"
	@echo "  make clean-logs       - Clear Docker logs"
	@echo "  make clean-images     - Remove Docker images"
	@echo "  make clean-all        - Full cleanup (DESTRUCTIVE!)"
	@echo "  make down-volumes     - Stop services and remove volumes (DESTRUCTIVE!)"
	@echo ""
	@echo "💾 Available Services:"
	@echo "  - gateway-service     - API Gateway (port 8080)"
	@echo "  - auth-service        - Authentication (port 8081)"
	@echo "  - stats-service       - System Statistics (port 8082)"
	@echo "  - postgres            - PostgreSQL Database (port 5432)"

.PHONY: init
init:
	@python3 $(PYTHON_SCRIPT)

# Start Docker services
.PHONY: up
up: init
	docker compose -f docker-compose.yml up -d

# Start Docker services with build (force rebuild)
.PHONY: up!
up!: build-ui init
	docker compose -f docker-compose.yml up -d --build

# Build all services
.PHONY: build
build: init
	docker compose -f docker-compose.yml build

# Show which services can be built vs use pre-built images
.PHONY: build-info
build-info:
	@echo "🔍 Service build information:"
	@echo "Services that can be built (have Dockerfiles):"
	@docker compose -f docker-compose.yml config | grep -B 1 "build:" | grep -E "^  [a-zA-Z-]+:" | sed 's/://g' | sed 's/^  /  ✅ /' || echo "  (none found)"
	@echo ""
	@echo "Services using pre-built images:"
	@docker compose -f docker-compose.yml config | grep -B 1 "image:" | grep -E "^  [a-zA-Z-]+:" | sed 's/://g' | sed 's/^  /  📦 /' || echo "  (none found)"

# Build specific service - Usage: make build-auth, make build-gateway
# Note: Only works for services with Dockerfiles (not postgres which uses pre-built image)
.PHONY: build-%
build-%: init
	@echo "Building $*..."
	@if docker compose -f docker-compose.yml config --services | grep -q "^$*$$"; then \
		if docker compose -f docker-compose.yml config | grep -A 10 "^  $*:" | grep -q "build:"; then \
			docker compose -f docker-compose.yml build $*; \
		else \
			echo "❌ Service '$*' uses a pre-built image (no Dockerfile to build)"; \
		fi \
	else \
		echo "❌ Service '$*' not found in docker-compose.yml"; \
	fi

# Build and start services (rebuild containers)
.PHONY: up-build
up-build: init
	docker compose -f docker-compose.yml up -d --build

# Stop Docker services (keeps volumes - data persists)
.PHONY: down
down:
	docker compose -f docker-compose.yml down

# Stop Docker services and remove volumes (DESTRUCTIVE!)
.PHONY: down-volumes
down-volumes:
	docker compose -f docker-compose.yml down -v

# Stop services gracefully (gives containers time to shutdown)
.PHONY: stop
stop:
	docker compose -f docker-compose.yml stop

# Start specific service - Usage: make postgres-up, make auth-up
.PHONY: %-up
%-up: init
	docker compose -f docker-compose.yml up -d $*

# Stop specific service - Usage: make postgres-down, make auth-down
.PHONY: %-down
%-down:
	docker compose -f docker-compose.yml stop $*
	docker compose -f docker-compose.yml rm -f $*

# Restart specific service - Usage: make postgres-restart, make auth-restart
.PHONY: %-restart
%-restart:
	docker compose -f docker-compose.yml restart $*

# View logs for specific service - Usage: make postgres-logs, make auth-logs
.PHONY: %-logs
%-logs:
	docker compose -f docker-compose.yml logs -f $*

# Login to Postgres container as root
.PHONY: postgres-login
postgres-login:
	docker exec -it postgres psql -U postgres

# Create a new user
.PHONY: create-user
create-user:
	@./scripts/create_user.sh

# List all users
.PHONY: list-users
list-users:
	@./scripts/list_users.sh

# Go module management
.PHONY: tidy
tidy:
	@echo "📦 Running go mod tidy for all Go services..."
	@for service in $(GO_SERVICES); do \
		if [ -f $$service/go.mod ]; then \
			echo "  ✅ Tidying $$service..."; \
			cd $$service && go mod tidy && cd ..; \
		fi \
	done
	@echo "✨ All Go modules tidied"

# Download Go dependencies for all services
.PHONY: deps
deps:
	@echo "📥 Downloading Go dependencies..."
	@for service in $(GO_SERVICES); do \
		if [ -f $$service/go.mod ]; then \
			echo "  ⬇️  Downloading $$service dependencies..."; \
			cd $$service && go mod download && cd ..; \
		fi \
	done
	@echo "✅ All dependencies downloaded"

# Verify Go modules
.PHONY: verify
verify:
	@echo "🔍 Verifying Go modules..."
	@for service in $(GO_SERVICES); do \
		if [ -f $$service/go.mod ]; then \
			echo "  🔎 Verifying $$service..."; \
			cd $$service && go mod verify && cd ..; \
		fi \
	done
	@echo "✅ All modules verified"

# Remove generated files
.PHONY: clean
clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf $(POSTGRES_INIT_DIR)
	rm -rf gateway/ui-build
	rm -f auth/auth-service
	rm -f gateway/gateway-service
	rm -f stats/stats-service
	@echo "✅ Cleanup complete"

# Clear all Docker images for services (DESTRUCTIVE!)
.PHONY: clean-images
clean-images:
	@echo "🗑️  Removing Docker images for services..."
	@docker compose -f docker-compose.yml down --rmi all 2>/dev/null || true
	@echo "✅ Docker images cleared"

# Clear all Docker logs
.PHONY: clean-logs
clean-logs:
	@echo "🧼 Clearing Docker logs..."
	rm -rf /tmp/home-server/*
	@echo "✅ Docker logs cleared"

# Full cleanup - stop services, remove volumes, images, and generated files (NUCLEAR!)
.PHONY: clean-all
clean-all: clean-images clean clean-logs
	@echo "🧹 Full cleanup completed"

# UI Build Commands
.PHONY: build-ui
build-ui:
	@echo "🎨 Building React UI..."
	@./scripts/build-ui.sh
