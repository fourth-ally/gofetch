# Contributing to GoFetch

Thank you for your interest in contributing to GoFetch!

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/fourth-ally/gofetch.git
   cd gofetch
   ```

2. **Install dependencies**
   ```bash
   make install
   ```

3. **Run tests**
   ```bash
   make test
   ```

## Code Quality Standards

### Test Coverage Requirement ⚠️

**All contributions must maintain minimum 80% test coverage.**

- Current coverage: **87.7%** ✅
- Total tests: 20

#### Running Tests

```bash
# Run all tests with coverage
make test

# Generate HTML coverage report
make coverage

# Check coverage percentage
go test -coverprofile=coverage.out -coverpkg=./infrastructure,./domain/... ./tests/...
go tool cover -func=coverage.out | tail -1
```

#### Writing Tests

- Place tests in the `tests/` directory
- Organize tests by feature category
- Use `httptest.Server` for HTTP mocking
- Ensure new features include corresponding tests
- Test both success and error cases

Example test structure:
```go
package tests

import (
    "testing"
    "github.com/fourth-ally/gofetch/infrastructure"
)

func TestMyFeature(t *testing.T) {
    // Arrange
    client := infrastructure.NewClient()
    
    // Act
    result := client.DoSomething()
    
    // Assert
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### Code Style

- Follow standard Go conventions
- Run `gofmt` before committing
- Use `go vet` to catch common issues
- Keep functions small and focused
- Add comments for exported functions

```bash
make fmt  # Format code
make vet  # Run static analysis
```

### Commit Messages

Use clear, descriptive commit messages:
- `feat: add new HTTP method support`
- `fix: correct path parameter handling`
- `test: add coverage for interceptors`
- `docs: update API documentation`

## Pull Request Process

1. **Fork the repository**

2. **Create a feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

3. **Make your changes**
   - Write code
   - Add tests (maintain 80% coverage)
   - Update documentation if needed

4. **Verify quality**
   ```bash
   make test      # All tests must pass
   make coverage  # Check coverage stays ≥80%
   make fmt       # Format code
   make vet       # Run static analysis
   ```

5. **Commit and push**
   ```bash
   git add .
   git commit -m "feat: add awesome feature"
   git push origin feature/my-feature
   ```

6. **Open a Pull Request**
   - Provide clear description of changes
   - Reference any related issues
   - Ensure CI checks pass

## Test Organization

Tests are organized by feature in the `tests/` directory:

- `common_test.go` - Shared test utilities
- `client_creation_test.go` - Client initialization
- `http_methods_test.go` - HTTP methods (GET, POST, PUT, PATCH, DELETE)
- `parameters_test.go` - Path and query parameters
- `error_handling_test.go` - Error scenarios
- `interceptors_test.go` - Request/response interceptors
- `context_test.go` - Context and cancellation
- `advanced_features_test.go` - Progress, transformers, config merge

## Coverage Breakdown

Current coverage by component:

| Component | Coverage |
|-----------|----------|
| Client configuration | 100% |
| HTTP methods | 100% |
| Interceptors | 100% |
| Progress tracking | 100% |
| Data transformers | 100% |
| URL building | 87.5% |
| Request execution | 80% |
| **Overall** | **87.7%** |

## Questions?

- Open an issue for bug reports or feature requests
- Check existing issues before creating new ones
- Be respectful and constructive

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
