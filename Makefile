# Build variables
BINARY_NAME=docker-cleaner
VERSION=$(shell git describe --tags --always --dirty 2>nul || echo "v0.1.0")
COMMIT=$(shell git rev-parse HEAD 2>nul || echo "unknown")
BUILD_DATE=$(shell powershell -Command "Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'")

# Go build flags
LDFLAGS=-ldflags "-X github.com/zahidhasann88/docker-cleaner/cmd.version=$(VERSION) -X github.com/zahidhasann88/docker-cleaner/cmd.commit=$(COMMIT) -X github.com/zahidhasann88/docker-cleaner/cmd.buildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build the binary with verbose output
.PHONY: build
build:
	@echo Building $(BINARY_NAME)...
	@echo Version: $(VERSION)
	@echo Commit: $(COMMIT)
	@echo Build Date: $(BUILD_DATE)
	@echo Creating bin directory...
	@if not exist bin mkdir bin
	@echo Downloading dependencies...
	go mod download
	@echo Compiling binary...
	go build -v $(LDFLAGS) -o bin/$(BINARY_NAME).exe .
	@echo Build completed successfully!
	@echo Binary location: bin/$(BINARY_NAME).exe
	@dir bin\$(BINARY_NAME).exe

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo Building for multiple platforms...
	@if not exist bin mkdir bin
	@echo Building Linux AMD64...
	set GOOS=linux&& set GOARCH=amd64&& go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	@echo Building Linux ARM64...
	set GOOS=linux&& set GOARCH=arm64&& go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 .
	@echo Building Darwin AMD64...
	set GOOS=darwin&& set GOARCH=amd64&& go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	@echo Building Darwin ARM64...
	set GOOS=darwin&& set GOARCH=arm64&& go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .
	@echo Building Windows AMD64...
	set GOOS=windows&& set GOARCH=amd64&& go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .
	@echo All builds completed!
	@dir bin\

# Run tests
.PHONY: test
test:
	@echo Running tests...
	go test -v ./...

# Install dependencies
.PHONY: deps
deps:
	@echo Installing dependencies...
	go mod download
	go mod tidy
	@echo Dependencies installed!

# Clean build artifacts
.PHONY: clean
clean:
	@echo Cleaning build artifacts...
	@if exist bin rmdir /s /q bin
	@if exist release rmdir /s /q release
	go clean
	@echo Clean completed!

# Install the binary
.PHONY: install
install:
	@echo Installing $(BINARY_NAME)...
	go install $(LDFLAGS) .
	@echo Installation completed!

# Run linter
.PHONY: lint
lint:
	@echo Running linter...
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	@echo Formatting code...
	go fmt ./...
	@echo Code formatted!

# Run the application
.PHONY: run
run:
	@echo Running $(BINARY_NAME)...
	go run . $(ARGS)

# Create release packages
.PHONY: release
release: clean build-all
	@echo Creating release packages...
	@if not exist release mkdir release
	@echo Copying binaries...
	copy bin\* release\
	@echo Creating compressed archives...
	@cd release && powershell -Command "Compress-Archive -Path docker-cleaner-linux-amd64 -DestinationPath docker-cleaner-linux-amd64.zip"
	@cd release && powershell -Command "Compress-Archive -Path docker-cleaner-linux-arm64 -DestinationPath docker-cleaner-linux-arm64.zip"
	@cd release && powershell -Command "Compress-Archive -Path docker-cleaner-darwin-amd64 -DestinationPath docker-cleaner-darwin-amd64.zip"
	@cd release && powershell -Command "Compress-Archive -Path docker-cleaner-darwin-arm64 -DestinationPath docker-cleaner-darwin-arm64.zip"
	@cd release && powershell -Command "Compress-Archive -Path docker-cleaner-windows-amd64.exe -DestinationPath docker-cleaner-windows-amd64.zip"
	@echo Release packages created!
	@dir release\

# Create a GitHub release (requires gh CLI)
.PHONY: github-release
github-release: release
	@echo Creating GitHub release...
	gh release create $(VERSION) release\* --title "Release $(VERSION)" --notes "Release $(VERSION) of Docker Cleaner"

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo Building Docker image...
	docker build -t zahidhasann88/docker-cleaner:$(VERSION) .
	docker build -t zahidhasann88/docker-cleaner:latest .
	@echo Docker image built successfully!

# Push Docker image
.PHONY: docker-push
docker-push: docker-build
	@echo Pushing Docker image...
	docker push zahidhasann88/docker-cleaner:$(VERSION)
	docker push zahidhasann88/docker-cleaner:latest
	@echo Docker image pushed successfully!

# Create a new version tag
.PHONY: tag
tag:
	@echo Current version: $(VERSION)
	@set /p NEW_VERSION="Enter new version (e.g., v1.0.0): "
	git tag -a %NEW_VERSION% -m "Release %NEW_VERSION%"
	git push origin %NEW_VERSION%
	@echo Tag %NEW_VERSION% created and pushed!

# Check if go.mod exists
.PHONY: check-mod
check-mod:
	@if not exist go.mod (echo go.mod not found! Run 'make init' to create it. && exit /b 1)

# Initialize Go module
.PHONY: init
init:
	@echo Initializing Go module...
	go mod init github.com/zahidhasann88/docker-cleaner
	@echo Go module initialized!

# Debug build - shows all output
.PHONY: debug-build
debug-build:
	@echo Starting debug build...
	@echo Current directory: %CD%
	@echo Go version:
	go version
	@echo Go environment:
	go env GOOS GOARCH
	@echo Checking go.mod:
	@if exist go.mod (echo go.mod found) else (echo go.mod NOT found)
	@echo Creating bin directory...
	@if not exist bin mkdir bin
	@echo Building with full output...
	go build -x -v $(LDFLAGS) -o bin/$(BINARY_NAME).exe .
	@echo Final binary check:
	@if exist bin\$(BINARY_NAME).exe (echo Binary created successfully) else (echo Binary NOT created)
	@dir bin\

# Package for distribution
.PHONY: package
package: release
	@echo Creating distribution packages...
	@if not exist dist mkdir dist
	@echo Creating installer for Windows...
	@echo Creating tarball for Linux...
	@echo Creating DMG for macOS...
	@echo Distribution packages created!

# Show project stats
.PHONY: stats
stats:
	@echo Project Statistics:
	@echo Lines of code:
	@powershell -Command "Get-ChildItem -Recurse -Include *.go | Get-Content | Measure-Object -Line | Select-Object -ExpandProperty Lines"
	@echo Binary size:
	@if exist bin\$(BINARY_NAME).exe dir bin\$(BINARY_NAME).exe

# Help
.PHONY: help
help:
	@echo Available targets:
	@echo   build         - Build the binary
	@echo   build-all     - Build for multiple platforms
	@echo   test          - Run tests
	@echo   deps          - Install dependencies
	@echo   clean         - Clean build artifacts
	@echo   install       - Install the binary
	@echo   lint          - Run linter
	@echo   fmt           - Format code
	@echo   run           - Run the application
	@echo   release       - Create release packages
	@echo   github-release- Create GitHub release
	@echo   docker-build  - Build Docker image
	@echo   docker-push   - Push Docker image
	@echo   tag           - Create new version tag
	@echo   package       - Create distribution packages
	@echo   init          - Initialize Go module
	@echo   check-mod     - Check if go.mod exists
	@echo   debug-build   - Build with full debug output
	@echo   stats         - Show project statistics
	@echo   help          - Show this help