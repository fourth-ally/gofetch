.PHONY: help build test example wasm wasm-serve clean fmt vet

help: ## Display this help message
	@echo "GoFetch - Makefile commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the library
	@echo "Building GoFetch..."
	@go build ./...
	@echo "✓ Build complete!"

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests complete!"

coverage: test ## Run tests with coverage report
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

example: ## Run the basic example
	@echo "Running basic example..."
	@go run examples/basic/main.go

wasm: ## Build WebAssembly binary
	@echo "Building for WebAssembly..."
	@./scripts/build-wasm.sh

wasm-serve: wasm ## Build and serve WASM demo
	@echo "Starting WASM demo server..."
	@cd examples/wasm && ./serve.sh

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted!"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet complete!"

lint: ## Run golangci-lint (requires golangci-lint installed)
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "✓ Lint complete!"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f coverage.out coverage.html
	@rm -rf dist/
	@rm -f examples/basic/main examples/basic/main.exe
	@echo "✓ Clean complete!"

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed!"

.DEFAULT_GOAL := help
