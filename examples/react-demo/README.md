# GoFetch React Demo

A demonstration of using GoFetch HTTP client in a React application via WebAssembly.

## Features

- ✨ Full GoFetch API in React components
- 🎯 Real-time request/response logging  
- 📊 Request statistics tracking
- 🎨 Modern UI with Vite + React
- 🚀 WebAssembly powered HTTP client

## Prerequisites

- Node.js 18+
- Go 1.21+ (for building WASM)

## Setup

1. **Install dependencies:**

```bash
npm install
```

2. **Build the GoFetch WASM module:**

```bash
npm run build:wasm
```

This will:
- Compile GoFetch to WebAssembly
- Copy the WASM binary to `public/gofetch.wasm` (~10MB)
- Copy the Go WASM runtime to `public/wasm_exec.js`

## Development

Start the development server:

```bash
npm run dev
```

Then open your browser to http://localhost:5173

## Usage Examples

The demo includes several interactive examples:

### Basic GET Request

```javascript
const response = await gofetch.get('/users')
console.log(response.data)
```

### GET with Query Parameters

```javascript
const response = await gofetch.get('/posts', { userId: 1 })
// Requests: /posts?userId=1
console.log(response.data)
```

### POST Request

```javascript
const newPost = {
  userId: 1,
  title: 'My Post',
  body: 'Post content'
}

const response = await gofetch.post('/posts', null, newPost)
console.log(response.data)
```

### Custom Client Instance

```javascript
const customClient = gofetch.newClient()
customClient.setBaseURL('https://api.example.com')
customClient.setTimeout(5000)
customClient.setHeader('Authorization', 'Bearer token')

const response = await customClient.get('/endpoint')
```

## Project Structure

```
react-demo/
├── public/              # Static assets
│   ├── gofetch.wasm    # GoFetch WebAssembly module
│   └── wasm_exec.js    # Go WASM runtime
├── src/
│   ├── hooks/
│   │   └── useGoFetch.js   # React hook for GoFetch
│   ├── App.jsx             # Main app component
│   ├── index.css           # Global styles
│   └── main.jsx            # React entry point
├── index.html
├── vite.config.js
└── package.json
```

## How It Works

1. **WASM Loading**: The `useGoFetch` hook loads the GoFetch WASM module on component mount
2. **Go Runtime**: `wasm_exec.js` provides the Go WebAssembly runtime
3. **API Exposure**: GoFetch exposes its API to JavaScript via `window.gofetch`
4. **Promise Wrapping**: All Go functions return JavaScript Promises
5. **React Integration**: The hook wraps the API with logging and state management

## API Reference

The `gofetch` object provides:

- `get(path, params)` - GET request
- `post(path, params, body)` - POST request
- `put(path, params, body)` - PUT request
- `patch(path, params, body)` - PATCH request
- `delete(path, params)` - DELETE request
- `setBaseURL(url)` - Configure base URL
- `setTimeout(ms)` - Set request timeout
- `setHeader(key, value)` - Set default header
- `newClient()` - Create new client instance

All methods return a Promise resolving to:

```javascript
{
  statusCode: 200,
  headers: { ... },
  data: { ... },
  rawBody: "..."
}
```

## Troubleshooting

### WASM Module Not Found

Run `npm run build:wasm` to build the WASM module.

### CORS Errors

The demo uses JSONPlaceholder API which has CORS enabled. For other APIs, ensure they have proper CORS headers.

## Building for Production

```bash
npm run build
npm run preview
```

The built files will be in the `dist` directory.

## License

MIT
