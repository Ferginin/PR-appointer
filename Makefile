.PHONY: help build run test test-coverage clean docker-build docker-up docker-down docker-logs swagger lint fmt deps

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application locally"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make lint           - Run linter (golangci-lint)"
	@echo "  make fmt            - Format code with gofmt"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start services with docker-compose"
	@echo "  make docker-down    - Stop services with docker-compose"
	@echo "  make docker-logs    - View docker-compose logs"
	@echo "  make migrate        - Run database migrations"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make deps           - Install dependencies"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/pr-appointer ./main.go
	@echo "Build complete: bin/pr-appointer"

# Run the application locally
run:
	@echo "Running application..."
	go run ./main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage report:"
	go tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML coverage report, run: go tool cover -html=coverage.out"

# Run linter
lint:
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run --config .golangci.yml ./...
	@echo "Linting complete"

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -w -s .
	goimports -w .
	@echo "Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out
	go clean
	@echo "Clean complete"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker-compose build
	@echo "Docker build complete"

# Start services with docker-compose
docker-up:
	@echo "Starting services..."
	docker-compose up -d
	@echo "Services started. Application available at http://localhost:8080"
	@echo "Health check: http://localhost:8080/health"

# Stop services
docker-down:
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped"

# View logs
docker-logs:
	docker-compose logs -f

# Restart services
docker-restart: docker-down docker-up

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@if ! command -v swag &> /dev/null; then \
		echo "swag not found. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag init -g main.go -o ./docs
	@echo "Swagger docs generated in ./docs"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

# Install dev tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Development tools installed"

# Full rebuild (clean + build)
rebuild: clean build

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	cp .env.example .env 2>/dev/null || echo ".env already exists"
	make deps
	make install-tools
	@echo "Development setup complete"
	@echo "Edit .env file with your configuration"

# Run all checks (format, lint, test)
check: fmt lint test
	@echo "All checks passed!"

# Pre-commit hook
pre-commit: fmt lint test
	@echo "Pre-commit checks passed!"
