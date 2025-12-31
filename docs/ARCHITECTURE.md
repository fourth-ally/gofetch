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
├── docs/                         # Documentation
│   ├── ARCHITECTURE.md          # This file - architecture overview
│   ├── CHANGELOG.md             # Version history and changes
│   ├── GO_USAGE.md              # Go usage examples
│   ├── NPM_PUBLISH.md           # NPM publishing guide
│   ├── PROJECT_SUMMARY.md       # Project overview
│   ├── QUICKSTART.md            # Quick start guide
│   └── RELEASE.md               # Release process
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
│   └── progress.go             # Progress tracking utilities
│
├── tests/                       # Test suite (organized by feature)
│   ├── client_creation_test.go # Client initialization tests
│   ├── http_methods_test.go    # GET, POST, etc. tests
│   ├── parameters_test.go      # Path & query parameter tests
│   ├── error_handling_test.go  # Error handling tests
│   ├── interceptors_test.go    # Interceptor tests
│   └── context_test.go         # Context & cancellation tests
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
│   ├── react-demo/             # React + Vite demo
│   │   ├── src/
│   │   │   └── hooks/
│   │   │       └── useGoFetch.js
│   │   └── package.json
│   ├── vue-demo/               # Vue + Vite demo
│   │   ├── src/
│   │   │   └── composables/
│   │   │       └── useGoFetch.js
│   │   └── package.json
│   └── wasm/                   # WebAssembly demo
│       ├── index.html          # Demo HTML page
│       └── serve.sh            # Local server script
│
└── scripts/                     # Build scripts
    ├── build-npm.js            # NPM package build script
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

**Coverage Requirement: Minimum 80%** ✅

The project maintains a strict test coverage policy:
- **Current Coverage**: 87.7% (exceeds 80% minimum)
- **Total Tests**: 20 comprehensive unit tests
- **Organization**: Tests organized by feature in separate files

### Test Categories

1. **Client Creation** (`client_creation_test.go`)
   - Client initialization and configuration
   - Fluent interface testing
   - Derived client instances

2. **HTTP Methods** (`http_methods_test.go`)
   - GET, POST, PUT, PATCH, DELETE requests
   - Request/response handling
   - JSON marshaling/unmarshaling

3. **Parameters** (`parameters_test.go`)
   - Path parameter substitution
   - Query parameter encoding
   - URL building

4. **Error Handling** (`error_handling_test.go`)
   - HTTP error responses
   - Custom status validators
   - Error message formatting

5. **Interceptors** (`interceptors_test.go`)
   - Request interceptor chain
   - Response interceptor chain
   - Interceptor execution order

6. **Context** (`context_test.go`)
   - Context cancellation
   - Timeout handling
   - Request abortion

7. **Advanced Features** (`advanced_features_test.go`)
   - Data transformers
   - Upload/download progress callbacks
   - Config merging
   - HTTPError.Error() method

### Testing Tools

- **HTTP Mocking**: Using `httptest.Server` for isolated testing
- **Coverage Analysis**: `go test -coverprofile=coverage.out`
- **Race Detection**: Tests run with `-race` flag
- **Continuous Coverage**: Coverage checked on every test run

### Coverage Policy

- **Minimum Required**: 80%
- **Current Level**: 87.7%
- **Enforcement**: Coverage reports generated with each test run
- **Focus Areas**: All public APIs must be tested

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
