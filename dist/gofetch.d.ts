// Type definitions for gofetch
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
