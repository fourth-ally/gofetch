# GoFetch - Quick Start Guide

Welcome to GoFetch! This guide will help you get started quickly.

## Installation

```bash
go get github.com/fourth-ally/gofetch
```

## 5-Minute Tutorial

### 1. Basic GET Request

```go
package main

import (
    "context"
    "fmt"
    "github.com/fourth-ally/gofetch"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func main() {
    client := gofetch.NewClient()
    
    var user User
    _, err := client.Get(
        context.Background(),
        "https://jsonplaceholder.typicode.com/users/1",
        nil,
        &user,
    )
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %s\n", user.Name)
}
```

### 2. Configure Base URL

```go
client := gofetch.NewClient().
    SetBaseURL("https://api.example.com").
    SetTimeout(10 * time.Second)

var users []User
_, err := client.Get(context.Background(), "/users", nil, &users)
```

### 3. POST Request

```go
newUser := User{Name: "John Doe"}

var createdUser User
_, err := client.Post(
    context.Background(),
    "/users",
    nil,          // no URL params
    newUser,      // request body
    &createdUser, // response target
)
```

### 4. URL Parameters

```go
// Path parameters
params := map[string]interface{}{"id": 123}
_, err := client.Get(ctx, "/users/:id", params, &user)

// Query parameters
params := map[string]interface{}{
    "page": 1,
    "limit": 10,
}
_, err := client.Get(ctx, "/users", params, &users)
// Requests: /users?page=1&limit=10
```

### 5. Error Handling

```go
var user User
_, err := client.Get(ctx, "/users/999", nil, &user)
if err != nil {
    if httpErr, ok := err.(*errors.HTTPError); ok {
        fmt.Printf("Status: %d\n", httpErr.StatusCode)
        fmt.Printf("Body: %s\n", string(httpErr.Body))
    }
}
```

### 6. Add Authentication

```go
client := gofetch.NewClient().
    SetBaseURL("https://api.example.com").
    SetHeader("Authorization", "Bearer YOUR_TOKEN")
```

### 7. Request Interceptor

```go
client := gofetch.NewClient().
    AddRequestInterceptor(func(req *http.Request) (*http.Request, error) {
        // Add custom header to every request
        req.Header.Set("X-Request-ID", generateID())
        return req, nil
    })
```

### 8. Progress Tracking

```go
client := gofetch.NewClient().
    SetDownloadProgress(func(transferred, total int64) {
        percent := float64(transferred) / float64(total) * 100
        fmt.Printf("\rProgress: %.2f%%", percent)
    })
```

## Common Patterns

### Creating Multiple Clients

```go
// Base client with common config
baseClient := gofetch.NewClient().
    SetBaseURL("https://api.example.com").
    SetTimeout(30 * time.Second)

// Authenticated client
authClient := baseClient.NewInstance().
    SetHeader("Authorization", "Bearer token")

// Admin client with different timeout
adminClient := baseClient.NewInstance().
    SetHeader("Authorization", "Bearer admin-token").
    SetTimeout(60 * time.Second)
```

### Handling Different Response Formats

```go
// If API wraps data in {"data": ...}
client := gofetch.NewClient().
    SetDataTransformer(func(body []byte) ([]byte, error) {
        var wrapper struct {
            Data json.RawMessage `json:"data"`
        }
        json.Unmarshal(body, &wrapper)
        return wrapper.Data, nil
    })
```

### Custom Success Criteria

```go
// Accept both 2xx and 3xx as success
client := gofetch.NewClient().
    SetStatusValidator(func(code int) bool {
        return code >= 200 && code < 400
    })
```

## Next Steps

1. Read the full [README.md](README.md) for comprehensive documentation
2. Explore [examples/basic/main.go](examples/basic/main.go) for more examples
3. Check out the [ARCHITECTURE.md](ARCHITECTURE.md) to understand the design
4. Try the WebAssembly demo in [examples/wasm/](examples/wasm/)

## Running Examples

```bash
# Run basic example
make example

# Or manually
go run examples/basic/main.go
```

## Building for WebAssembly

```bash
# Build WASM
make wasm

# Serve demo
make wasm-serve
# Then open http://localhost:8080
```

## Getting Help

- Check the examples directory for common use cases
- Read the comprehensive README
- Review the test files for usage patterns

## Quick Reference

| Method | Purpose |
|--------|---------|
| `NewClient()` | Create new client |
| `SetBaseURL(url)` | Set base URL |
| `SetTimeout(duration)` | Set timeout |
| `SetHeader(k, v)` | Set default header |
| `Get(ctx, path, params, target)` | GET request |
| `Post(ctx, path, params, body, target)` | POST request |
| `Put(ctx, path, params, body, target)` | PUT request |
| `Patch(ctx, path, params, body, target)` | PATCH request |
| `Delete(ctx, path, params, target)` | DELETE request |
| `NewInstance()` | Create derived client |

Happy coding with GoFetch! 🚀
