package infrastructure

import (
	"sync"
	"time"

	"github.com/fourth-ally/gofetch/domain/models"
)

// CircuitBreaker implements circuit breaker pattern for HTTP endpoints.
type CircuitBreaker struct {
	mu sync.RWMutex

	// Circuit state per endpoint (URL)
	circuits map[string]*circuitState

	// Configuration
	threshold        int
	timeout          time.Duration
	halfOpenRequests int
}

// circuitState tracks the state of a single circuit.
type circuitState struct {
	state                models.CircuitBreakerState
	failureCount         int
	lastFailureTime      time.Time
	openedAt             time.Time
	halfOpenSuccessCount int
	halfOpenAttempts     int
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(threshold int, timeout time.Duration, halfOpenRequests int) *CircuitBreaker {
	return &CircuitBreaker{
		circuits:         make(map[string]*circuitState),
		threshold:        threshold,
		timeout:          timeout,
		halfOpenRequests: halfOpenRequests,
	}
}

// IsOpen checks if the circuit is open for a given endpoint.
func (cb *CircuitBreaker) IsOpen(endpoint string) bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	circuit, exists := cb.circuits[endpoint]
	if !exists {
		return false
	}

	// Check if circuit should transition from open to half-open
	if circuit.state == models.CircuitBreakerOpen {
		if time.Since(circuit.openedAt) >= cb.timeout {
			return false // Allow transition to half-open
		}
		return true
	}

	return false
}

// CanAttempt checks if a request can be attempted (respects half-open limit).
func (cb *CircuitBreaker) CanAttempt(endpoint string) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	circuit, exists := cb.circuits[endpoint]
	if !exists {
		// Create new circuit in closed state
		cb.circuits[endpoint] = &circuitState{
			state: models.CircuitBreakerClosed,
		}
		return true
	}

	// Transition from open to half-open if timeout expired
	if circuit.state == models.CircuitBreakerOpen {
		if time.Since(circuit.openedAt) >= cb.timeout {
			circuit.state = models.CircuitBreakerHalfOpen
			circuit.halfOpenAttempts = 0
			circuit.halfOpenSuccessCount = 0
		} else {
			return false
		}
	}

	// In half-open state, limit concurrent requests
	if circuit.state == models.CircuitBreakerHalfOpen {
		if circuit.halfOpenAttempts >= cb.halfOpenRequests {
			return false
		}
		circuit.halfOpenAttempts++
	}

	return true
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess(endpoint string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	circuit, exists := cb.circuits[endpoint]
	if !exists {
		return
	}

	if circuit.state == models.CircuitBreakerHalfOpen {
		circuit.halfOpenSuccessCount++
		// If we got enough successes, close the circuit
		if circuit.halfOpenSuccessCount >= cb.halfOpenRequests {
			circuit.state = models.CircuitBreakerClosed
			circuit.failureCount = 0
		}
	} else if circuit.state == models.CircuitBreakerClosed {
		// Reset failure count on success
		circuit.failureCount = 0
	}
}

// RecordFailure records a failed request.
func (cb *CircuitBreaker) RecordFailure(endpoint string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	circuit, exists := cb.circuits[endpoint]
	if !exists {
		circuit = &circuitState{
			state: models.CircuitBreakerClosed,
		}
		cb.circuits[endpoint] = circuit
	}

	circuit.failureCount++
	circuit.lastFailureTime = time.Now()

	// Open circuit if threshold exceeded
	if circuit.state == models.CircuitBreakerClosed && circuit.failureCount >= cb.threshold {
		circuit.state = models.CircuitBreakerOpen
		circuit.openedAt = time.Now()
	}

	// If half-open request failed, reopen the circuit
	if circuit.state == models.CircuitBreakerHalfOpen {
		circuit.state = models.CircuitBreakerOpen
		circuit.openedAt = time.Now()
		circuit.halfOpenAttempts = 0
		circuit.halfOpenSuccessCount = 0
	}
}

// GetState returns the current state of a circuit.
func (cb *CircuitBreaker) GetState(endpoint string) models.CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	circuit, exists := cb.circuits[endpoint]
	if !exists {
		return models.CircuitBreakerClosed
	}

	return circuit.state
}

// Reset resets the circuit breaker for a specific endpoint.
func (cb *CircuitBreaker) Reset(endpoint string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	delete(cb.circuits, endpoint)
}

// ResetAll resets all circuit breakers.
func (cb *CircuitBreaker) ResetAll() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.circuits = make(map[string]*circuitState)
}
