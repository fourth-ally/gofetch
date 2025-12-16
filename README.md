# GoFetch

A robust, developer-friendly Go HTTP client library inspired by Axios, providing high-level configuration, automated data handling, and superior error management.

## ✨ Features

- **Fluent Configuration API**: Chain methods for intuitive client setup
- **Automatic JSON Handling**: Automatic marshaling and unmarshaling of request/response data
- **Request & Response Interceptors**: Modify requests and responses before/after execution
- **Data Transformers**: Transform response data before unmarshaling
- **URL Parameter Handling**: Support for both path variables (`:id`) and query parameters
- **Superior Error Management**: Custom error types with full response details
- **Progress Tracking**: Track upload and download progress
- **Context Integration**: Full support for cancellation and timeouts
- **WebAssembly Compatible**: Run GoFetch in the browser with WASM
- **Domain-Driven Design**: Clean architecture with separated concerns

## 📦 Installation

```bash
go get github.com/nikoskechris/gofetch
```

## 🚀 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/nikoskechris/gofetch"
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

## 📖 Usage Examples

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

## 🌐 WebAssembly Support

GoFetch can be compiled to WebAssembly and used in browsers:

### Building for WASM

```bash
GOOS=js GOARCH=wasm go build -o gofetch.wasm ./cmd/wasm
```

### Using in JavaScript

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <script src="wasm_exec.js"></script>
</head>
<body>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(
            fetch("gofetch.wasm"),
            go.importObject
        ).then((result) => {
            go.run(result.instance);
            
            // Use GoFetch
            gofetch.setBaseURL("https://api.example.com");
            
            gofetch.get("/users/1")
                .then(response => {
                    console.log("User:", response.data);
                })
                .catch(error => {
                    console.error("Error:", error);
                });
        });
    </script>
</body>
</html>
```

### Creating WASM Client Instances

```javascript
// Create a new client instance
const client = gofetch.newClient();

client.setBaseURL("https://api.example.com");
client.setTimeout(10000); // 10 seconds
client.setHeader("Authorization", "Bearer token123");

// Make requests
client.get("/users/1")
    .then(response => console.log(response))
    .catch(error => console.error(error));

client.post("/users", null, { name: "John", email: "john@example.com" })
    .then(response => console.log("Created:", response.data))
    .catch(error => console.error(error));
```

## 🏗️ Architecture

GoFetch follows Domain-Driven Design principles with a clean separation of concerns:

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

## 📋 API Reference

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

## 🧪 Testing

Run the example:

```bash
go run examples/basic/main.go
```

Build for WASM:

```bash
GOOS=js GOARCH=wasm go build -o gofetch.wasm ./cmd/wasm
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 🙏 Acknowledgments

- Inspired by [Axios](https://axios-http.com/) - the beloved HTTP client for JavaScript
- Built with Go's excellent standard library

## 📚 Additional Resources

- [Examples Directory](./examples/) - More usage examples
- [Domain Documentation](./domain/) - Understanding the domain layer
- [WASM Guide](./wasm/) - Detailed WebAssembly integration guide
