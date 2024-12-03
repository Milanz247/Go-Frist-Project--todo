# Define project-related variables
BINARY_NAME=myapp
MAIN_FILE=cmd/main.go
BUILD_DIR=bin
GO_FILES=$(shell find . -type f -name '*.go')
DSN="user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"

# Default target
.PHONY: all
all: build

# Live reload with Go Air
.PHONY: live
live:
	@echo "Starting live reload with Air..."
	~/go/bin/air -c air.toml
# Build the project
.PHONY: build
build:
	@echo "Building project..."
	set GOOS=linux&& set GOARCH=amd64&& go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete! Binary located at $(BUILD_DIR)/$(BINARY_NAME)"

# Test the project
.PHONY: test
test:
	@echo "Running tests..."
	go test ./... -cover
	@echo "Tests complete!"

# Clean the build
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Cleanup complete!"

# Lint the project
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Linting complete!"

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Formatting complete!"

# Watch for changes and trigger tests
.PHONY: watch
watch:
	@echo "Watching for changes to run tests..."
	reflex -r '\.go$$' -- sh -c 'go test ./...'

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed!"

# Help message
.PHONY: help
help:
	@echo "Available make targets:"
	@echo "  make           - Default target, builds the project"
	@echo "  make live      - Start live reloading with Air"
	@echo "  make build     - Build the Go project"
	@echo "  make test      - Run Go tests"
	@echo "  make clean     - Clean the build output"
	@echo "  make lint      - Run linter (using golangci-lint)"
	@echo "  make fmt       - Format Go code"
	@echo "  make watch     - Watch for file changes and run tests"
	@echo "  make deps      - Install dependencies"
	@echo "  make db-reset  - Drop all tables and reset the database"
	@echo "  make db-migrate - Migrate the database schema"
	@echo "  make db-seed   - Seed the database with initial data"
	@echo "  make help      - Show this help message"

# Declare file dependencies for auto-rebuilds (useful if a file changes)
$(BINARY_NAME): $(GO_FILES)
	@$(MAKE) build

# Database reset: Drop all tables
.PHONY: db-reset
db-reset:
	@echo "Resetting the database (dropping all tables)..."
	go run db/migrations.go --reset
	@echo "Database reset complete!"

# Database migration: AutoMigrate Go models
.PHONY: db-migrate
db-migrate:
	@echo "Migrating the database schema..."
	go run db/migrations.go --migrate
	@echo "Database migration complete!"

# Seed the database with initial data
.PHONY: db-seed
db-seed:
	@echo "Seeding the database..."
	go run db/migrations.go --seed
	@echo "Database seeding complete!"
