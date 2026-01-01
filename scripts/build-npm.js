const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

console.log('🔨 Building GoFetch for npm...\n');

// Create dist directory
const distDir = path.join(__dirname, '..', 'dist');
if (!fs.existsSync(distDir)) {
  fs.mkdirSync(distDir, { recursive: true });
}

// Build WASM
console.log('📦 Compiling Go to WebAssembly...');
try {
  execSync('GOOS=js GOARCH=wasm go build -o dist/gofetch.wasm ./cmd/wasm', {
    stdio: 'inherit',
    cwd: path.join(__dirname, '..')
  });
  console.log('✅ WASM compilation complete\n');
} catch (error) {
  console.error('❌ WASM compilation failed');
  process.exit(1);
}

// Copy wasm_exec.js
console.log('📋 Copying wasm_exec.js...');
let wasmExecContent = '';
try {
  // Try multiple possible locations
  const possiblePaths = [
    '/usr/local/go/misc/wasm/wasm_exec.js',
    '/opt/homebrew/opt/go/libexec/misc/wasm/wasm_exec.js'
  ];
  
  // Also try GOROOT
  try {
    const goRoot = execSync('go env GOROOT', { encoding: 'utf-8' }).trim();
    possiblePaths.unshift(path.join(goRoot, 'misc', 'wasm', 'wasm_exec.js'));
  } catch (e) {
    // Ignore if go env fails
  }
  
  let copied = false;
  
  for (const wasmExecSrc of possiblePaths) {
    if (fs.existsSync(wasmExecSrc)) {
      wasmExecContent = fs.readFileSync(wasmExecSrc, 'utf-8');
      // We'll inline this content into the JS wrapper instead of copying
      console.log('✅ wasm_exec.js content loaded\n');
      copied = true;
      break;
    }
  }
  
  if (!copied) {
    console.log('⚠️  Could not find wasm_exec.js, checked:');
    possiblePaths.forEach(p => console.log(`   - ${p}`));
    throw new Error('wasm_exec.js not found');
  }
} catch (error) {
  console.error('❌ Failed to load wasm_exec.js:', error.message);
  process.exit(1);
}

// Create JavaScript wrapper
console.log('📝 Creating JavaScript wrapper...');
const jsWrapper = `// GoFetch - HTTP Client Library
// Compiled from Go to WebAssembly

// Inline wasm_exec.js content
${wasmExecContent}

let gofetchInstance = null;
let initPromise = null;

async function initGoFetch() {
  if (gofetchInstance) {
    return gofetchInstance;
  }

  if (initPromise) {
    return initPromise;
  }

  initPromise = (async () => {
    console.log('GoFetch: Starting initialization...');
    
    if (typeof Go === 'undefined') {
      throw new Error('Go runtime not available - wasm_exec.js failed to load');
    }
    console.log('GoFetch: Go runtime available');

    const go = new Go();
    console.log('GoFetch: Go instance created');
    
    // Fetch WASM file
    const wasmUrl = new URL('./gofetch.wasm', import.meta.url);
    console.log('GoFetch: Fetching WASM from', wasmUrl.href);
    const response = await fetch(wasmUrl);
    const wasmBuffer = await response.arrayBuffer();
    console.log('GoFetch: WASM loaded, size:', wasmBuffer.byteLength);

    const result = await WebAssembly.instantiate(wasmBuffer, go.importObject);
    console.log('GoFetch: WASM instantiated');
    
    // Run the Go program
    go.run(result.instance);
    console.log('GoFetch: Go program started');

    // Wait for the gofetch global to be available
    const maxAttempts = 200; // 10 seconds timeout
    let attempts = 0;
    
    while (attempts < maxAttempts) {
      const globalScope = typeof window !== 'undefined' ? window : 
                         typeof global !== 'undefined' ? global : 
                         typeof globalThis !== 'undefined' ? globalThis : null;
      
      if (globalScope && globalScope.gofetch) {
        console.log('GoFetch: Found gofetch global after', attempts, 'attempts');
        gofetchInstance = globalScope.gofetch;
        break;
      }
      
      await new Promise(resolve => setTimeout(resolve, 50));
      attempts++;
      
      if (attempts % 20 === 0) {
        console.log('GoFetch: Still waiting for initialization... attempt', attempts);
      }
    }

    if (!gofetchInstance) {
      console.error('GoFetch: Initialization failed after', attempts, 'attempts');
      throw new Error('GoFetch WASM module failed to initialize after 10 seconds');
    }

    console.log('GoFetch: Initialization complete!');
    return gofetchInstance;
  })();

  return initPromise;
}

// Export API
export async function newClient() {
  const gf = await initGoFetch();
  return gf.newClient();
}

export async function get(url, params) {
  const gf = await initGoFetch();
  return gf.get(url, params);
}

export async function post(url, params, body) {
  const gf = await initGoFetch();
  return gf.post(url, params, body);
}

export async function put(url, params, body) {
  const gf = await initGoFetch();
  return gf.put(url, params, body);
}

export async function patch(url, params, body) {
  const gf = await initGoFetch();
  return gf.patch(url, params, body);
}

export async function del(url, params) {
  const gf = await initGoFetch();
  return gf.delete(url, params);
}

export async function setBaseURL(url) {
  const gf = await initGoFetch();
  return gf.setBaseURL(url);
}

export async function setTimeout(ms) {
  const gf = await initGoFetch();
  return gf.setTimeout(ms);
}

export async function setHeader(key, value) {
  const gf = await initGoFetch();
  return gf.setHeader(key, value);
}

export async function setRetryOptions(options) {
  const gf = await initGoFetch();
  return gf.setRetryOptions(options);
}

// Default export
export default {
  newClient,
  get,
  post,
  put,
  patch,
  delete: del,
  setBaseURL,
  setTimeout,
  setHeader,
  setRetryOptions
};
`;

fs.writeFileSync(path.join(distDir, 'gofetch.js'), jsWrapper);
console.log('✅ JavaScript wrapper created\n');

// Create TypeScript definitions
console.log('📝 Creating TypeScript definitions...');
const tsDefinitions = `// Type definitions for gofetch
// Project: https://github.com/fourth-ally/gofetch

export interface GoFetchResponse {
  statusCode: number;
  headers: Record<string, string | string[]>;
  data: any;
  rawBody: string;
}

export type BackoffStrategy = 'exponential' | 'linear' | 'fixed';

export interface RetryOptions {
  /** Maximum number of retry attempts (0 = no retries) */
  maxRetries?: number;
  /** Initial delay before first retry in milliseconds */
  initialDelay?: number;
  /** Maximum delay between retries in milliseconds */
  maxDelay?: number;
  /** Backoff strategy: exponential, linear, or fixed */
  backoff?: BackoffStrategy;
  /** Enable random jitter to prevent thundering herd */
  jitter?: boolean;
  /** Fraction of delay to randomize (0.0 - 1.0, default 0.3) */
  jitterFraction?: number;
  /** Additional HTTP status codes to retry (e.g., [429, 503]) */
  retryOnStatusCodes?: number[];
  /** Enable circuit breaker functionality */
  circuitBreaker?: boolean;
  /** Number of consecutive failures before opening circuit */
  circuitBreakerThreshold?: number;
  /** How long circuit stays open in milliseconds */
  circuitBreakerTimeout?: number;
  /** Number of requests allowed in half-open state */
  circuitBreakerHalfOpenRequests?: number;
}

export interface GoFetchClient {
  get(path: string, params?: Record<string, any>): Promise<GoFetchResponse>;
  post(path: string, params?: Record<string, any>, body?: any): Promise<GoFetchResponse>;
  put(path: string, params?: Record<string, any>, body?: any): Promise<GoFetchResponse>;
  patch(path: string, params?: Record<string, any>, body?: any): Promise<GoFetchResponse>;
  delete(path: string, params?: Record<string, any>): Promise<GoFetchResponse>;
  setBaseURL(url: string): GoFetchClient;
  setTimeout(ms: number): GoFetchClient;
  setHeader(key: string, value: string): GoFetchClient;
  setRetryOptions(options: RetryOptions): GoFetchClient;
  newInstance(): GoFetchClient;
}

export function newClient(): Promise<GoFetchClient>;
export function get(url: string, params?: Record<string, any>): Promise<GoFetchResponse>;
export function post(url: string, params?: Record<string, any>, body?: any): Promise<GoFetchResponse>;
export function put(url: string, params?: Record<string, any>, body?: any): Promise<GoFetchResponse>;
export function patch(url: string, params?: Record<string, any>, body?: any): Promise<GoFetchResponse>;
export function del(url: string, params?: Record<string, any>): Promise<GoFetchResponse>;
export function setBaseURL(url: string): Promise<void>;
export function setTimeout(ms: number): Promise<void>;
export function setHeader(key: string, value: string): Promise<void>;

declare const gofetch: {
  newClient: typeof newClient;
  get: typeof get;
  post: typeof post;
  put: typeof put;
  patch: typeof patch;
  delete: typeof del;
  setBaseURL: typeof setBaseURL;
  setTimeout: typeof setTimeout;
  setHeader: typeof setHeader;
};

export default gofetch;
`;

fs.writeFileSync(path.join(distDir, 'gofetch.d.ts'), tsDefinitions);
console.log('✅ TypeScript definitions created\n');

console.log('🎉 Build complete! Package ready for npm publish.\n');
console.log('📦 To publish:');
console.log('   npm login');
console.log('   npm publish\n');
