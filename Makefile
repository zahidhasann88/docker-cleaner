# Build variables
BINARY_NAME=docker-cleaner
VERSION=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build flags
LDFLAGS=-ldflags "-X github.com/yourusername/docker-cleaner/cmd.version=$(VERSION) -X github.com/yourusername/docker-cleaner/cmd.commit=$(COMMIT) -X github.com/yourusername/docker-cleaner/cmd.buildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .

# Run tests
.PHONY: test
test:
	go test -v ./...

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	go clean

# Install the binary
.PHONY: install
install:
	go install $(LDFLAGS) .

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run the application
.PHONY: run
run:
	go run . $(ARGS)

# Create release
.PHONY: release
release: clean build-all
	mkdir -p release
	cp bin/* release/
	cd release && for f in *; do tar -czf $$f.tar.gz $$f; done
	cd release && rm docker-cleaner-*
	ls -la release/

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  test       - Run tests"
	@echo "  deps       - Install dependencies"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install the binary"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  run        - Run the application"
	@echo "  release    - Create release packages"
	@echo "  help       - Show this help"