# GoFetch - Project Structure

This document provides an overview of the GoFetch project structure following Domain-Driven Design principles.

## Architecture Overview

GoFetch follows a clean, layered architecture inspired by Domain-Driven Design:

```
┌─────────────────────────────────────────────────┐
│           Application Layer                     │
│  (gofetch.go - Public API Entry Point)         │
└─────────────────────────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│         Infrastructure Layer                    │
│  (HTTP Client Implementation)                   │
│  - Client, Request Execution                    │
│  - Progress Tracking                            │
└─────────────────────────────────────────────────┘
                     ▼
┌─────────────────────────────────────────────────┐
│            Domain Layer                         │
│  - Models (Config, Response)                    │
│  - Contracts (Interfaces)                       │
│  - Errors (HTTPError)                           │
└─────────────────────────────────────────────────┘
```

## Directory Structure

```
gofetch/
│
├── gofetch.go                    # Public API entry point
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
├── README.md                     # Project documentation
├── LICENSE                       # MIT license
├── Makefile                      # Build automation
├── .gitignore                    # Git ignore rules
│
├── domain/                       # Domain layer (core business logic)
│   ├── models/                   # Domain models
│   │   ├── config.go            # Configuration model
│   │   └── response.go          # Response model
│   ├── contracts/               # Interfaces and contracts
│   │   └── interceptors.go      # Interceptor contracts
│   └── errors/                  # Domain errors
│       └── http_error.go        # HTTP error type
│
├── infrastructure/              # Infrastructure layer (implementation)
│   ├── client.go               # HTTP client implementation
│   ├── client_test.go          # Client tests
│   └── progress.go             # Progress tracking utilities
│
├── wasm/                        # WebAssembly bridge
│   ├── bridge.go               # JavaScript bridge functions
│   └── helpers.go              # WASM utility functions
│
├── cmd/                         # Command-line applications
│   └── wasm/                   # WASM build entry point
│       └── main.go             # WASM main function
│
├── examples/                    # Usage examples
│   ├── basic/                  # Basic usage example
│   │   └── main.go
│   └── wasm/                   # WebAssembly demo
│       ├── index.html          # Demo HTML page
│       └── serve.sh            # Local server script
│
└── scripts/                     # Build scripts
    └── build-wasm.sh           # WASM build script
```

## Layer Responsibilities

### Domain Layer (`domain/`)

The domain layer contains the core business logic and is independent of external frameworks or libraries.

**Models** (`domain/models/`)
- `Config`: Client configuration with merge capabilities
- `Response`: HTTP response wrapper

**Contracts** (`domain/contracts/`)
- `RequestInterceptor`: Request modification interface
- `ResponseInterceptor`: Response inspection interface
- `DataTransformer`: Response data transformation interface
- `ProgressCallback`: Progress tracking interface

**Errors** (`domain/errors/`)
- `HTTPError`: Rich HTTP error type with status code, body, and headers

### Infrastructure Layer (`infrastructure/`)

The infrastructure layer implements the domain contracts and provides the actual HTTP client functionality.

**Client** (`infrastructure/client.go`)
- HTTP client implementation
- Fluent configuration API
- Request execution with interceptors
- URL building (path params & query strings)
- JSON marshaling/unmarshaling
- Error handling

**Progress** (`infrastructure/progress.go`)
- Progress tracking for uploads/downloads
- `progressReader` implementation

### WebAssembly Bridge (`wasm/`)

The WebAssembly bridge exposes GoFetch to JavaScript environments.

**Bridge** (`wasm/bridge.go`)
- JavaScript function exposure
- Promise wrapping for async operations
- Client instance management

**Helpers** (`wasm/helpers.go`)
- JavaScript ↔ Go type conversion
- Promise wrapper implementation
- Response transformation

### Application Layer

**Public API** (`gofetch.go`)
- Single entry point: `NewClient()`
- Clean, simple public interface

## Key Design Patterns

### 1. Fluent Interface (Builder Pattern)

```go
client := gofetch.NewClient().
    SetBaseURL("https://api.example.com").
    SetTimeout(10 * time.Second).
    SetHeader("Authorization", "Bearer token")
```

### 2. Chain of Responsibility (Interceptors)

```go
client.AddRequestInterceptor(authInterceptor).
       AddRequestInterceptor(loggingInterceptor).
       AddResponseInterceptor(metricsInterceptor)
```

### 3. Strategy Pattern (StatusValidator)

```go
client.SetStatusValidator(func(statusCode int) bool {
    return statusCode >= 200 && statusCode < 400
})
```

### 4. Template Method (Request Execution)

The `executeRequest` method defines the algorithm skeleton:
1. Merge configuration
2. Build URL
3. Marshal request body
4. Apply request interceptors
5. Execute HTTP request
6. Apply response interceptors
7. Validate status
8. Transform data
9. Unmarshal response

### 5. Factory Pattern (NewClient, NewInstance)

```go
baseClient := gofetch.NewClient()
derivedClient := baseClient.NewInstance()
```

## Testing Strategy

- **Unit Tests**: Comprehensive tests for all client functionality
- **HTTP Mocking**: Using `httptest.Server` for isolated testing
- **Test Coverage**: All major paths including error cases

## Build Commands

```bash
# Build the library
make build

# Run tests
make test

# Run with coverage
make coverage

# Build WebAssembly
make wasm

# Serve WASM demo
make wasm-serve

# Format code
make fmt

# Run static analysis
make vet
```

## WebAssembly Compilation

GoFetch can be compiled to WebAssembly:

```bash
GOOS=js GOARCH=wasm go build -o gofetch.wasm ./cmd/wasm
```

The build uses the `// +build js,wasm` constraint to include WASM-specific code.

## Dependencies

- **Zero External Dependencies**: Uses only Go standard library
- **Minimal Surface Area**: Clean, focused API
- **WebAssembly Compatible**: All code works in WASM environment

## Extension Points

The architecture makes it easy to extend GoFetch:

1. **Custom Interceptors**: Implement authentication, logging, metrics
2. **Custom Transformers**: Parse specialized API response formats
3. **Custom Status Validators**: Define success conditions per API
4. **Progress Callbacks**: Track large file uploads/downloads

## Performance Considerations

- **Connection Reuse**: Uses `http.Client` with connection pooling
- **Context Integration**: Supports cancellation and timeouts
- **Streaming**: Progress tracking with minimal memory overhead
- **Zero Allocations**: Efficient path parameter replacement

## Security Best Practices

- **Context Timeouts**: Prevents hanging requests
- **Status Validation**: Automatic error detection
- **Header Control**: Full control over request headers
- **Body Inspection**: Access to raw response bodies for validation
