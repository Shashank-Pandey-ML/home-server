# Paths
POSTGRES_INIT_DIR=postgres
PYTHON_SCRIPT=generate_postgres_setup.py

# Go services
GO_SERVICES=auth gateway stats common

# Default target
.PHONY: all
all: build-ui up!

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
	rm -rf $(POSTGRES_INIT_DIR)
	rm -rf gateway/ui-build

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
