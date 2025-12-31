package models

import "time"

// BackoffStrategy defines the strategy for calculating retry delays.
type BackoffStrategy string

const (
	// BackoffExponential doubles the delay after each retry.
	BackoffExponential BackoffStrategy = "exponential"
	// BackoffLinear increases the delay linearly after each retry.
	BackoffLinear BackoffStrategy = "linear"
	// BackoffFixed uses a fixed delay between retries.
	BackoffFixed BackoffStrategy = "fixed"
)

// RetryOptions configures the retry behavior for HTTP requests.
type RetryOptions struct {
	// MaxRetries is the maximum number of retry attempts (0 = no retries).
	MaxRetries int

	// InitialDelay is the initial delay before the first retry.
	InitialDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration

	// Backoff is the strategy used to calculate retry delays.
	Backoff BackoffStrategy

	// Jitter enables random jitter to prevent thundering herd.
	Jitter bool

	// JitterFraction is the fraction of delay to randomize (0.0 - 1.0).
	// Default is 0.3 (30% of the delay).
	JitterFraction float64

	// RetryOnStatusCodes specifies additional status codes to retry.
	// By default, only 5xx errors are retried.
	RetryOnStatusCodes []int

	// CircuitBreaker enables circuit breaker functionality.
	CircuitBreaker bool

	// CircuitBreakerThreshold is the number of consecutive failures
	// before opening the circuit.
	CircuitBreakerThreshold int

	// CircuitBreakerTimeout is how long the circuit stays open.
	CircuitBreakerTimeout time.Duration

	// CircuitBreakerHalfOpenRequests is the number of requests to allow
	// in half-open state to test if the service recovered.
	CircuitBreakerHalfOpenRequests int
}

// NewRetryOptions creates default retry options.
func NewRetryOptions() *RetryOptions {
	return &RetryOptions{
		MaxRetries:                     3,
		InitialDelay:                   100 * time.Millisecond,
		MaxDelay:                       30 * time.Second,
		Backoff:                        BackoffExponential,
		Jitter:                         true,
		JitterFraction:                 0.3,
		RetryOnStatusCodes:             []int{},
		CircuitBreaker:                 false,
		CircuitBreakerThreshold:        5,
		CircuitBreakerTimeout:          60 * time.Second,
		CircuitBreakerHalfOpenRequests: 1,
	}
}

// ShouldRetryStatus checks if a status code should trigger a retry.
func (r *RetryOptions) ShouldRetryStatus(statusCode int) bool {
	// Always retry 5xx errors
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// Check custom retry status codes
	for _, code := range r.RetryOnStatusCodes {
		if statusCode == code {
			return true
		}
	}

	return false
}

// CircuitBreakerState represents the state of a circuit breaker.
type CircuitBreakerState string

const (
	// CircuitBreakerClosed allows requests through normally.
	CircuitBreakerClosed CircuitBreakerState = "closed"
	// CircuitBreakerOpen blocks requests and returns errors immediately.
	CircuitBreakerOpen CircuitBreakerState = "open"
	// CircuitBreakerHalfOpen allows limited requests to test recovery.
	CircuitBreakerHalfOpen CircuitBreakerState = "half-open"
)
