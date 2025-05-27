# MCP Registry CLI Makefile

.PHONY: help build clean install test lint fmt vet deps

# Default target
help:
	@echo "Available targets:"
	@echo "  build    - Build the CLI binary"
	@echo "  clean    - Remove build artifacts"
	@echo "  install  - Install dependencies"
	@echo "  test     - Run tests"
	@echo "  lint     - Run linter"
	@echo "  fmt      - Format code"
	@echo "  vet      - Run go vet"
	@echo "  deps     - Download dependencies"

# Build the CLI binary
build:
	go build -o bin/mcp-cli .

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
install: deps
	go install .

# Run tests
test:
	go test ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Download dependencies
deps:
	go mod download
	go mod tidy
