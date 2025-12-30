# gofetch-wasm

A robust HTTP client library compiled from Go to WebAssembly for use in JavaScript/TypeScript applications.

## Installation

```bash
npm install gofetch-wasm
```

## Disclaimer
This library is experimental,beta version and provided as-is.
No guarantees are made regarding stability or correctness.
Use at your own risk.

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
We suggest using Vite for the best developer experience and compatibility with WebAssembly. See the Vite configuration note below for a short code snippet to add to your `vite.config.ts`.

### Vite configuration (important)

If your project uses Vite, the dev server may try to process the WASM file and serve it incorrectly (HTML instead of binary). To ensure the WebAssembly file is handled properly, add the following settings to your `vite.config.ts`:

```ts
// vite.config.ts
export default defineConfig({
  // ...your existing config
  optimizeDeps: {
    exclude: ['gofetch-wasm'],
  },
  assetsInclude: ['**/*.wasm'],
});
```

If you still see a `WebAssembly.instantiate(): expected magic word` error in the browser, make sure the WASM file is served with the `application/wasm` content-type and consider copying `node_modules/gofetch-wasm/dist/gofetch.wasm` into your `public/` directory or adding a small dev-server middleware to set the header during development.


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

## Browser Usage & Performance Considerations

`gofetch-wasm` runs a Go-based HTTP client inside the browser via WebAssembly.  
This provides consistent HTTP behavior across environments, but comes with important trade-offs.

### Bundle Size & Runtime Overhead

- Includes a WebAssembly binary and Go runtime
- Adds additional bundle size compared to native `fetch` or Axios
- Introduces JS ↔ WASM bridge overhead per request

Because of this, `gofetch-wasm` **is not a lightweight replacement for `fetch`**.

### When to Use in the Browser

Recommended for:
- Internal tools
- Admin dashboards
- B2B applications
- Developer tools
- SDKs that require identical HTTP behavior across browser, Node.js, and edge runtimes

In these cases, the added overhead is usually acceptable and the consistency benefits outweigh the cost.

### When *Not* to Use in the Browser

Not recommended for:
- Consumer-facing applications
- Landing pages or marketing sites
- Mobile-first apps
- Performance-critical UIs
- Applications with high-frequency or real-time HTTP requests

For these use cases, native `fetch` or Axios will provide better performance and smaller bundles.

### Design Philosophy

`gofetch-wasm` prioritizes **deterministic behavior and cross-runtime consistency** over minimal bundle size or raw performance.  
It is designed as an engineering trade-off, not a general-purpose browser HTTP client.

## License

MIT

