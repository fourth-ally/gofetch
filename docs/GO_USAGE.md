# GoFetch - Go Usage Guide

Complete usage guide for using GoFetch as a Go library.

## Installation

```bash
go get github.com/fourth-ally/gofetch
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/fourth-ally/gofetch"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Create a new client
    client := gofetch.NewClient().
        SetBaseURL("https://api.example.com").
        SetTimeout(10 * time.Second).
        SetHeader("Authorization", "Bearer token123")
    
    // Make a GET request
    var user User
    resp, err := client.Get(context.Background(), "/users/1", nil, &user)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
}
```

## Usage Examples

### Basic GET Request

```go
client := gofetch.NewClient().
    SetBaseURL("https://api.example.com")

var users []User
resp, err := client.Get(context.Background(), "/users", nil, &users)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Retrieved %d users\n", len(users))
```

### POST Request with Body

```go
newUser := User{
    Name:  "John Doe",
    Email: "john@example.com",
}

var createdUser User
resp, err := client.Post(context.Background(), "/users", nil, newUser, &createdUser)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created user with ID: %d\n", createdUser.ID)
```

### Path Parameters

```go
params := map[string]interface{}{
    "id": 123,
}

var user User
resp, err := client.Get(context.Background(), "/users/:id", params, &user)
```

### Query Parameters

```go
params := map[string]interface{}{
    "page":     1,
    "per_page": 10,
    "status":   "active",
}

var users []User
resp, err := client.Get(context.Background(), "/users", params, &users)
// Request URL: /users?page=1&per_page=10&status=active
```

### Request Interceptors

```go
client := gofetch.NewClient().
    AddRequestInterceptor(func(req *http.Request) (*http.Request, error) {
        // Add custom header
        req.Header.Set("X-Request-ID", generateRequestID())
        
        // Log request
        log.Printf("Making request to %s", req.URL.String())
        
        return req, nil
    })
```

### Response Interceptors

```go
client := gofetch.NewClient().
    AddResponseInterceptor(func(resp *http.Response) (*http.Response, error) {
        // Log response
        log.Printf("Received response with status %d", resp.StatusCode)
        
        // Check custom headers
        if rateLimitRemaining := resp.Header.Get("X-RateLimit-Remaining"); rateLimitRemaining != "" {
            log.Printf("Rate limit remaining: %s", rateLimitRemaining)
        }
        
        return resp, nil
    })
```

### Data Transformers

```go
// Transform response data to extract payload from wrapper
client := gofetch.NewClient().
    SetDataTransformer(func(data []byte) ([]byte, error) {
        var wrapper struct {
            Data json.RawMessage `json:"data"`
        }
        
        if err := json.Unmarshal(data, &wrapper); err != nil {
            return data, nil // Return original if not wrapped
        }
        
        return wrapper.Data, nil
    })

// Now API responses like {"data": {...}} will automatically unwrap
```

### Error Handling

```go
var user User
_, err := client.Get(context.Background(), "/users/99999", nil, &user)
if err != nil {
    if httpErr, ok := err.(*errors.HTTPError); ok {
        fmt.Printf("HTTP Error: Status %d\n", httpErr.StatusCode)
        fmt.Printf("Response body: %s\n", string(httpErr.Body))
        fmt.Printf("Headers: %v\n", httpErr.Headers)
    } else {
        fmt.Printf("Request error: %v\n", err)
    }
}
```

### Custom Status Validation

```go
// Accept 2xx and 3xx as success
client := gofetch.NewClient().
    SetStatusValidator(func(statusCode int) bool {
        return statusCode >= 200 && statusCode < 400
    })
```

### Progress Tracking

```go
client := gofetch.NewClient().
    SetDownloadProgress(func(bytesTransferred, totalBytes int64) {
        percentage := float64(bytesTransferred) / float64(totalBytes) * 100
        fmt.Printf("\rDownload progress: %.2f%%", percentage)
    }).
    SetUploadProgress(func(bytesTransferred, totalBytes int64) {
        percentage := float64(bytesTransferred) / float64(totalBytes) * 100
        fmt.Printf("\rUpload progress: %.2f%%", percentage)
    })
```

### Creating Derived Clients

```go
// Base client with common configuration
baseClient := gofetch.NewClient().
    SetBaseURL("https://api.example.com").
    SetHeader("X-App-Version", "1.0.0").
    SetTimeout(30 * time.Second)

// Create authenticated client
authClient := baseClient.NewInstance().
    SetHeader("Authorization", "Bearer token123")

// Create admin client with longer timeout
adminClient := baseClient.NewInstance().
    SetHeader("Authorization", "Bearer admin-token").
    SetTimeout(60 * time.Second)
```

### Context and Cancellation

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var users []User
_, err := client.Get(ctx, "/users", nil, &users)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(2 * time.Second)
    cancel() // Cancel request after 2 seconds
}()

_, err := client.Get(ctx, "/users", nil, &users)
```

## Architecture

GoFetch uses a layered architecture with a clear separation of concerns:

```
gofetch/
├── domain/              # Domain layer - pure business logic
│   ├── contracts/       # Interfaces and contracts
│   ├── errors/          # Domain error types
│   └── models/          # Domain models
├── infrastructure/      # Infrastructure layer - implementations
│   ├── client.go        # HTTP client implementation
│   └── progress.go      # Progress tracking utilities
├── wasm/                # WebAssembly bridge
│   ├── bridge.go        # JavaScript bridge
│   └── helpers.go       # WASM utilities
├── examples/            # Usage examples
└── gofetch.go          # Public API entry point
```

## API Reference

### Client Methods

#### Configuration

- `NewClient() *Client` - Create a new client instance
- `SetBaseURL(url string) *Client` - Set base URL for all requests
- `SetTimeout(duration time.Duration) *Client` - Set request timeout
- `SetHeader(key, value string) *Client` - Set default header
- `SetStatusValidator(func(int) bool) *Client` - Set custom status validator
- `NewInstance() *Client` - Create derived client with inherited settings

#### Interceptors & Transformers

- `AddRequestInterceptor(RequestInterceptor) *Client` - Add request interceptor
- `AddResponseInterceptor(ResponseInterceptor) *Client` - Add response interceptor
- `SetDataTransformer(DataTransformer) *Client` - Set data transformer

#### Progress Tracking

- `SetUploadProgress(ProgressCallback) *Client` - Set upload progress callback
- `SetDownloadProgress(ProgressCallback) *Client` - Set download progress callback

#### HTTP Methods

- `Get(ctx, path, params, target) (*Response, error)` - Perform GET request
- `Post(ctx, path, params, body, target) (*Response, error)` - Perform POST request
- `Put(ctx, path, params, body, target) (*Response, error)` - Perform PUT request
- `Patch(ctx, path, params, body, target) (*Response, error)` - Perform PATCH request
- `Delete(ctx, path, params, target) (*Response, error)` - Perform DELETE request

### Types

```go
type Response struct {
    StatusCode int
    Headers    http.Header
    Data       interface{}
    RawBody    []byte
}

type HTTPError struct {
    StatusCode   int
    Body         []byte
    Headers      http.Header
    Message      string
    OriginalResp *http.Response
}

type RequestInterceptor func(*http.Request) (*http.Request, error)
type ResponseInterceptor func(*http.Response) (*http.Response, error)
type DataTransformer func([]byte) ([]byte, error)
type ProgressCallback func(bytesTransferred, totalBytes int64)
```

## Building for WebAssembly

Build the WASM binary:

```bash
GOOS=js GOARCH=wasm go build -o gofetch.wasm ./cmd/wasm
```

## Testing

Run the example:

```bash
go run examples/basic/main.go
```
