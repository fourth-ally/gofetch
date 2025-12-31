# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [1.0.12] - TBD

### Added
- **Retry Logic**: Automatic request retry with configurable max attempts
  - Three backoff strategies: exponential, linear, and fixed
  - Configurable initial delay and max delay
  - Optional jitter to prevent thundering herd (default 30%)
  - Retry on 5xx errors by default
  - Custom status codes for retry (e.g., 429 Too Many Requests)
- **Circuit Breaker**: Per-endpoint failure tracking to prevent cascading failures
  - Configurable failure threshold before opening circuit
  - Automatic transition to half-open state after timeout
  - Configurable number of requests allowed in half-open state
  - Independent circuit tracking for each endpoint
  - Works with or without retry logic (can use circuit breaker alone)
- **New Domain Models**: `RetryOptions`, `BackoffStrategy`, `CircuitBreakerState`
- **New Infrastructure Components**: `RetryManager`, `CircuitBreaker`
- **WASM Support**: Full JavaScript/TypeScript API for retry and circuit breaker
- **Comprehensive Tests**: 11 new tests covering all retry and circuit breaker scenarios

### Changed
- Client now supports `SetRetryOptions()` method for configuring retry behavior
- All HTTP methods now use retry wrapper when configured
- Test coverage maintained at 80.8% (above 80% minimum requirement)
- Total test count increased from 20 to 31 tests

### Fixed
- Circuit breaker now works independently when `MaxRetries=0`
- Circuit breaker checks execute before retry logic to prevent unnecessary attempts
