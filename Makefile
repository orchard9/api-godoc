# API GoDoc - Makefile
# Build automation for OpenAPI documentation generator

# Binary name and version
BINARY_NAME=api-godoc
VERSION?=dev
BUILD_TIME=$(shell date -u +%Y%m%d.%H%M%S)
BUILD_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.buildHash=$(BUILD_HASH)"

# Directories
BUILD_DIR=build
CMD_DIR=cmd/$(BINARY_NAME)
COVERAGE_DIR=coverage

.PHONY: all build clean test coverage lint fmt vet deps tidy ci uat help ci-setup act-test release-local

# Default target
all: ci

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run linters
lint: fmt vet
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping additional linting"; \
	fi

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...

# Tidy modules
tidy:
	@echo "Tidying modules..."
	$(GOMOD) tidy

# CI target - comprehensive validation
ci: deps tidy lint test build
	@echo "CI pipeline completed successfully"

# User Acceptance Testing
uat: build
	@echo "Running User Acceptance Testing..."
	@if [ ! -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "Error: Binary not found. Run 'make build' first."; \
		exit 1; \
	fi
	@echo "Testing --help flag..."
	./$(BUILD_DIR)/$(BINARY_NAME) --help
	@echo "Testing --version flag..."
	./$(BUILD_DIR)/$(BINARY_NAME) --version
	@echo "Testing with UAT artifacts..."
	@if [ -f uat/artifacts/warden.v1.swagger.json ]; then \
		echo "Testing with warden.v1.swagger.json..."; \
		./$(BUILD_DIR)/$(BINARY_NAME) uat/artifacts/warden.v1.swagger.json; \
	fi
	@if [ -f uat/artifacts/forge.swagger.json ]; then \
		echo "Testing with forge.swagger.json..."; \
		./$(BUILD_DIR)/$(BINARY_NAME) uat/artifacts/forge.swagger.json; \
	fi
	@echo "UAT completed successfully"

# Run comprehensive UAT tests
.PHONY: uat-test
uat-test:
	@cd uat && go run runner.go

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin; \
	fi

# Quick development build and test
dev: fmt vet test build
	@echo "Development build completed"

# Install the binary to GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/

# CI setup target
ci-setup:
	@echo "Setting up CI environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin; \
	fi

# Test GitHub Actions locally with ACT
act-test:
	@echo "Testing GitHub Actions locally with ACT..."
	@if ! command -v act >/dev/null 2>&1; then \
		echo "Error: ACT not installed. Install it first: curl -q https://raw.githubusercontent.com/nektos/act/master/install.sh | bash"; \
		exit 1; \
	fi
	@echo "Testing CI workflow..."
	./bin/act -j test --platform ubuntu-latest=catthehacker/ubuntu:act-latest
	@echo "Testing lint workflow..."
	./bin/act -j lint --platform ubuntu-latest=catthehacker/ubuntu:act-latest

# Build release binaries locally
release-local:
	@echo "Building release binaries locally..."
	@mkdir -p dist
	@echo "Building Linux amd64..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	@echo "Building Linux arm64..."
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)
	@echo "Building macOS amd64..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	@echo "Building macOS arm64..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 ./$(CMD_DIR)
	@echo "Building Windows amd64..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "Release binaries built in dist/"

# Help target
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  lint         - Run linters"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  deps         - Install dependencies"
	@echo "  tidy         - Tidy modules"
	@echo "  ci           - Run full CI pipeline"
	@echo "  uat          - Run User Acceptance Testing"
	@echo "  dev          - Quick development build"
	@echo "  dev-setup    - Set up development environment"
	@echo "  ci-setup     - Set up CI environment"
	@echo "  act-test     - Test GitHub Actions locally with ACT"
	@echo "  release-local - Build release binaries locally"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  help         - Show this help message"