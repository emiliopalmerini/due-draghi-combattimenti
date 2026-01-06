.PHONY: help build run test clean templ fmt vet lint dev install-tools all

BINARY_NAME=combattimenti
BINARY_PATH=bin/$(BINARY_NAME)
CMD_PATH=cmd/encounters/main.go
PORT?=8080

help:
	@echo "Available targets:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application (default port 8080)"
	@echo "  make dev            - Run in development mode (templ generate + run)"
	@echo "  make test           - Run tests"
	@echo "  make templ          - Generate templ templates"
	@echo "  make fmt            - Format Go code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run all linters (fmt + vet)"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make install-tools  - Install required tools (templ)"
	@echo "  make all            - Clean, generate templates, and build"

build: fmt vet templ
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_PATH) $(CMD_PATH)
	@echo "Build complete: $(BINARY_PATH)"

run: build
	@echo "Starting $(BINARY_NAME) on port $(PORT)..."
	@./$(BINARY_PATH)

dev:
	@echo "Running in development mode..."
	@templ generate --watch --cmd="go run $(CMD_PATH)"

test: fmt vet
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

test-short:
	@echo "Running short tests..."
	@go test -v -short ./...

templ:
	@echo "Generating templ templates..."
	@templ generate

fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

lint: fmt vet
	@echo "Linting complete"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Clean complete"

install-tools:
	@echo "Installing templ..."
	@go install github.com/a-h/templ/cmd/templ@latest
	@echo "Tools installed"

all: clean build
	@echo "All tasks complete"
