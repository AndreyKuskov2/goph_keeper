BINARY_NAME=goph_keeper_client
VERSION=1.0.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"

.PHONY: all
all: build

.PHONY: build
build:
@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/goph_keeper_client

.PHONY: build-all
build-all: build-windows build-linux build-darwin

.PHONY: build-windows
build-windows:
@echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_windows_amd64.exe ./cmd/goph_keeper_client
GOOS=windows GOARCH=386 go build $(LDFLAGS) -o $(BINARY_NAME)_windows_386.exe ./cmd/goph_keeper_client

.PHONY: build-linux
build-linux:
@echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 ./cmd/goph_keeper_client
GOOS=linux GOARCH=386 go build $(LDFLAGS) -o $(BINARY_NAME)_linux_386 ./cmd/goph_keeper_client
GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)_linux_arm64 ./cmd/goph_keeper_client

.PHONY: build-darwin
build-darwin:
@echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)_darwin_amd64 ./cmd/goph_keeper_client
GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)_darwin_arm64 ./cmd/goph_keeper_client

.PHONY: clean
clean:
@echo "Cleaning build artifacts..."
rm -f $(BINARY_NAME)
rm -f $(BINARY_NAME)_*

.PHONY: deps
deps:
@echo "Installing dependencies..."
go mod tidy
go mod download

.PHONY: test
test:
@echo "Running tests..."
go test ./...
