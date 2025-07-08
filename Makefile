# Build variables
BINARY_NAME=docker-cleaner
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "unknown")

# Go build flags
LDFLAGS=-ldflags "-X github.com/zahidhasann88/docker-cleaner/cmd.version=$(VERSION) -X github.com/zahidhasann88/docker-cleaner/cmd.commit=$(COMMIT) -X github.com/zahidhasann88/docker-cleaner/cmd.buildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"
	@if [ ! -d "bin" ]; then mkdir bin; fi
	go mod download
	go build -v $(LDFLAGS) -o bin/$(BINARY_NAME).exe .
	@echo "Build completed successfully!"
	@echo "Binary location: bin/$(BINARY_NAME).exe"

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@if [ ! -d "bin" ]; then mkdir bin; fi
	@echo "Building Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	@echo "Building Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 .
	@echo "Building Darwin AMD64..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	@echo "Building Darwin ARM64..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .
	@echo "Building Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .
	@echo "All builds completed!"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@if [ -d "bin" ]; then rm -rf bin; fi
	@if [ -d "release" ]; then rm -rf release; fi
	go clean
	@echo "Clean completed!"

# Install the binary
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .
	@echo "Installation completed!"

# Run linter
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	go run . $(ARGS)

# Create release packages
.PHONY: release
release: clean build-all
	@echo "Creating release packages..."
	@if [ ! -d "release" ]; then mkdir release; fi
	@echo "Copying binaries..."
	@cp bin/* release/ 2>/dev/null || true
	@echo "Release packages created!"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  test          - Run tests"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  install       - Install the binary"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  run           - Run the application"
	@echo "  release       - Create release packages"
	@echo "  help          - Show this help"