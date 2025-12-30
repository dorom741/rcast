# ==============================================================================
# rcast - RCast DMR
# ==============================================================================

# Project Information
PROJECT_NAME := rcast
BINARY_NAME := $(PROJECT_NAME)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build Configuration
GO := go
GO_VERSION := $(shell $(GO) version | cut -d' ' -f3)
GO_PACKAGES := $(shell $(GO) list ./...)
MAIN_PACKAGE := .

# Directories
OUTPUT_DIR := output
BIN_DIR := bin
LOG_DIR := $(OUTPUT_DIR)/log

# Build Flags
LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.buildTime=$(BUILD_TIME) \
	-X main.gitCommit=$(GIT_COMMIT) \
	-X main.goVersion=$(GO_VERSION)

GOBUILD_FLAGS :=-buildvcs=false -trimpath  -ldflags="$(LDFLAGS)"

# Development Configuration
GO_ENV ?= production

# Default target
.DEFAULT_GOAL := help

# ==============================================================================
# Build Targets
# ==============================================================================

# Build all platforms
.PHONY: build-all
build-all: build-linux-arm64 build-linux-amd64 build-windows-amd64

# Build for Linux ARM64 (musl)
.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "Building for Linux ARM64..."
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(GOBUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 main.go

# Build for Linux AMD64 (musl)
.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "Building for Linux AMD64..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GOBUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 main.go

# Build for Windows AMD64
.PHONY: build-windows-amd64
build-windows-amd64:
	@echo "Building for Windows AMD64..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(GOBUILD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go


.PHONY: build-dev
build-dev: prepare ## Build development binary (with debug info)
	@echo "Building development version of $(BINARY_NAME)..."
	@GOWORK=off $(GO) build -o $(BIN_DIR)/$(BINARY_NAME)-dev $(MAIN_PACKAGE)
	@echo "Development build completed: $(BIN_DIR)/$(BINARY_NAME)-dev"

# ==============================================================================
# Development Targets
# ==============================================================================

.PHONY: run
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@$(BIN_DIR)/$(BINARY_NAME)

.PHONY: run-dev
run-dev: build-dev ## Build and run development version
	@echo "Running development version..."
	@$(BIN_DIR)/$(BINARY_NAME)-dev --debug

# ==============================================================================
# Testing & Quality Targets
# ==============================================================================

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	@$(GO) test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@$(GO) test -race -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@golangci-lint run

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...

.PHONY: fmt
fmt: ## Format Go code using goimports-reviser
	@echo "Formatting Go code with goimports-reviser..."
	@goimports-reviser -format ./...

.PHONY: tidy
tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	@$(GO) mod tidy

# ==============================================================================
# Utility Targets
# ==============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)/bin/*
	@rm -f coverage.out coverage.html
	@echo "Clean completed"

.PHONY: clean-all
clean-all: ## Clean all output directories
	@echo "Cleaning all output..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Full clean completed"

.PHONY: prepare
prepare: ## Prepare output directories
	@mkdir -p $(BIN_DIR) $(LOG_DIR)

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GO) mod download

.PHONY: version
version: ## Show version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"

# ==============================================================================
# Help Target
# ==============================================================================

.PHONY: help
help: ## Show this help message
	@echo "$(PROJECT_NAME) - RCast DMR"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ==============================================================================
# Development Workflow Shortcuts
# ==============================================================================

.PHONY: dev
dev: tidy fmt vet lint test build ## Complete development workflow
	@echo "Development workflow completed successfully"

.PHONY: ci
ci: tidy fmt vet lint test-coverage build ## CI workflow
	@echo "CI workflow completed successfully"

# ==============================================================================
# End of Makefile
# ==============================================================================
