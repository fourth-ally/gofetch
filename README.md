# gofetch-wasm

A robust HTTP client library compiled from Go to WebAssembly for use in JavaScript/TypeScript applications.

## Installation

```bash
npm install gofetch-wasm
```

## Quick Start

```javascript
import gofetch from 'gofetch-wasm'

// Initialize the client
const client = await gofetch.newClient()

// Configure the client
client.setBaseURL('https://api.example.com')
client.setTimeout(10000) // 10 seconds
client.setHeader('Authorization', 'Bearer token123')

// Make requests
const response = await client.get('/users/1')
console.log(response.data)
```

## API

### Client Configuration

```javascript
const client = await gofetch.newClient()

client.setBaseURL(url)              // Set base URL
client.setTimeout(milliseconds)     // Set timeout
client.setHeader(key, value)        // Set default header
```

### HTTP Methods

```javascript
// GET request
const response = await client.get(path, params)

// POST request
const response = await client.post(path, params, body)

// PUT request
const response = await client.put(path, params, body)

// PATCH request
const response = await client.patch(path, params, body)

// DELETE request
const response = await client.delete(path, params)
```

### Parameters

```javascript
// Query parameters
const params = { page: 1, limit: 10 }
await client.get('/users', params)
// GET /users?page=1&limit=10

// Path parameters
const params = { id: 123 }
await client.get('/users/:id', params)
// GET /users/123

// Request body
const body = { name: 'John', email: 'john@example.com' }
await client.post('/users', null, body)
```

### Response Format

```javascript
{
  statusCode: 200,
  headers: { ... },
  data: { ... }  // Parsed JSON response
}
```

### Error Handling

```javascript
try {
  const response = await client.get('/users/999')
} catch (error) {
  console.error('Request failed:', error)
}
```

## TypeScript Support

TypeScript definitions are included:

```typescript
import gofetch from 'gofetch-wasm'

interface User {
  id: number
  name: string
  email: string
}

const client = await gofetch.newClient()
const response = await client.get('/users/1')
const user = response.data as User
```

## License

MIT
