// Type definitions for gofetch
// Project: https://github.com/nkechr1s/gofetch

export interface GoFetchResponse {
  statusCode: number;
  headers: Record<string, string | string[]>;
  data: any;
  rawBody: string;
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
