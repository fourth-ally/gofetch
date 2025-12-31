package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fourth-ally/gofetch/domain/models"
	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestRetryBasicSuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetRetryOptions(&models.RetryOptions{
			MaxRetries:   3,
			InitialDelay: 10 * time.Millisecond,
			MaxDelay:     100 * time.Millisecond,
			Backoff:      models.BackoffExponential,
			Jitter:       false,
		})

	var result map[string]string
	resp, err := client.Get(context.Background(), "/test", nil, &result)

	if err != nil {
		t.Fatalf("Expected request to succeed after retries, got error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	if result["status"] != "success" {
		t.Errorf("Expected success status, got %v", result)
	}
}

func TestRetryExhaustion(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetRetryOptions(&models.RetryOptions{
			MaxRetries:   2,
			InitialDelay: 10 * time.Millisecond,
			MaxDelay:     100 * time.Millisecond,
			Backoff:      models.BackoffFixed,
			Jitter:       false,
		})

	_, err := client.Get(context.Background(), "/test", nil, nil)

	if err == nil {
		t.Fatal("Expected error after retry exhaustion")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts (initial + 2 retries), got %d", attempts)
	}
}

func TestNoRetryWithoutConfiguration(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)

	_, err := client.Get(context.Background(), "/test", nil, nil)

	if err == nil {
		t.Fatal("Expected error")
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retries without config), got %d", attempts)
	}
}

func TestExponentialBackoff(t *testing.T) {
	retryManager := infrastructure.NewRetryManager(&models.RetryOptions{
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     10 * time.Second,
		Backoff:      models.BackoffExponential,
		Jitter:       false,
	})

	delays := []time.Duration{
		retryManager.CalculateDelay(0),
		retryManager.CalculateDelay(1),
		retryManager.CalculateDelay(2),
		retryManager.CalculateDelay(3),
	}

	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}

	for i, delay := range delays {
		if delay != expected[i] {
			t.Errorf("Attempt %d: expected delay %v, got %v", i, expected[i], delay)
		}
	}
}

func TestCircuitBreakerOpens(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetRetryOptions(&models.RetryOptions{
			MaxRetries:              0,
			CircuitBreaker:          true,
			CircuitBreakerThreshold: 3,
			CircuitBreakerTimeout:   1 * time.Second,
		})

	ctx := context.Background()

	// Make threshold number of requests to trigger circuit breaker
	for i := 0; i < 3; i++ {
		_, err := client.Get(ctx, "/test", nil, nil)
		if err != nil && strings.Contains(err.Error(), "circuit breaker is open") {
			t.Errorf("Circuit opened too early at request %d, threshold is 3", i+1)
			return
		}
	}

	// Next request should be blocked by circuit breaker
	_, err := client.Get(ctx, "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected circuit breaker to block request")
	}
	if !strings.Contains(err.Error(), "circuit breaker is open") {
		t.Errorf("Expected circuit breaker error, got: %v", err)
	}
}

func TestCircuitBreakerPerEndpoint(t *testing.T) {
	attempts := make(map[string]int)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts[r.URL.Path]++

		if r.URL.Path == "/fail" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetRetryOptions(&models.RetryOptions{
			MaxRetries:              0,
			CircuitBreaker:          true,
			CircuitBreakerThreshold: 2,
			CircuitBreakerTimeout:   1 * time.Second,
		})

	ctx := context.Background()

	// Fail /fail endpoint enough times to open its circuit
	for i := 0; i < 2; i++ {
		client.Get(ctx, "/fail", nil, nil)
	}

	// Next request to /fail should be blocked
	_, err := client.Get(ctx, "/fail", nil, nil)
	if err == nil || !strings.Contains(err.Error(), "circuit breaker is open") {
		t.Error("Expected /fail circuit to be open")
	}

	// /success should still work
	_, err = client.Get(ctx, "/success", nil, nil)
	if err != nil {
		t.Errorf("Expected /success to work, got error: %v", err)
	}

	if attempts["/success"] != 1 {
		t.Errorf("Expected 1 request to /success, got %d", attempts["/success"])
	}
}

func TestRetryOnCustomStatusCodes(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests) // 429
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetRetryOptions(&models.RetryOptions{
			MaxRetries:         3,
			InitialDelay:       10 * time.Millisecond,
			RetryOnStatusCodes: []int{429},
		})

	resp, err := client.Get(context.Background(), "/test", nil, nil)

	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

// TestLinearBackoff tests the linear backoff calculation
func TestLinearBackoff(t *testing.T) {
	retryOpts := models.NewRetryOptions()
	retryOpts.InitialDelay = 100 * time.Millisecond
	retryOpts.MaxDelay = 10 * time.Second
	retryOpts.Backoff = models.BackoffLinear
	retryOpts.Jitter = false

	manager := infrastructure.NewRetryManager(retryOpts)

	tests := []struct {
		attempt  int
		expected time.Duration
	}{
		{0, 100 * time.Millisecond},
		{1, 200 * time.Millisecond},
		{2, 300 * time.Millisecond},
		{3, 400 * time.Millisecond},
	}

	for _, tt := range tests {
		delay := manager.CalculateDelay(tt.attempt)
		if delay != tt.expected {
			t.Errorf("Attempt %d: expected delay %v, got %v", tt.attempt, tt.expected, delay)
		}
	}
}

// TestFixedBackoff tests the fixed backoff calculation
func TestFixedBackoff(t *testing.T) {
	retryOpts := models.NewRetryOptions()
	retryOpts.InitialDelay = 100 * time.Millisecond
	retryOpts.MaxDelay = 10 * time.Second
	retryOpts.Backoff = models.BackoffFixed
	retryOpts.Jitter = false

	manager := infrastructure.NewRetryManager(retryOpts)

	for attempt := 0; attempt < 5; attempt++ {
		delay := manager.CalculateDelay(attempt)
		if delay != 100*time.Millisecond {
			t.Errorf("Attempt %d: expected fixed delay 100ms, got %v", attempt, delay)
		}
	}
}

// TestJitter tests that jitter adds randomization to delays
func TestJitter(t *testing.T) {
	retryOpts := models.NewRetryOptions()
	retryOpts.InitialDelay = 1000 * time.Millisecond
	retryOpts.MaxDelay = 10 * time.Second
	retryOpts.Backoff = models.BackoffFixed
	retryOpts.Jitter = true
	retryOpts.JitterFraction = 0.3

	manager := infrastructure.NewRetryManager(retryOpts)

	// Test multiple times to ensure jitter is working
	delays := make(map[time.Duration]bool)
	for i := 0; i < 10; i++ {
		delay := manager.CalculateDelay(0)
		delays[delay] = true

		// With 30% jitter, delay should be between 700ms and 1300ms
		if delay < 700*time.Millisecond || delay > 1300*time.Millisecond {
			t.Errorf("Delay %v is outside expected jitter range [700ms, 1300ms]", delay)
		}
	}

	// Should have at least some variation (not all the same)
	if len(delays) < 2 {
		t.Error("Expected variation in delays with jitter enabled")
	}
}

// TestCircuitBreakerRecovery tests that circuit breaker transitions to half-open and recovers
func TestCircuitBreakerRecovery(t *testing.T) {
	failUntil := time.Now().Add(400 * time.Millisecond)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Before(failUntil) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "recovered"}`))
	}))
	defer server.Close()

	client := infrastructure.NewClient()
	client.SetBaseURL(server.URL)

	client.SetRetryOptions(&models.RetryOptions{
		MaxRetries:                     0,
		CircuitBreaker:                 true,
		CircuitBreakerThreshold:        2,
		CircuitBreakerTimeout:          300 * time.Millisecond,
		CircuitBreakerHalfOpenRequests: 1,
	})

	ctx := context.Background()

	// Trip the circuit breaker
	for i := 0; i < 2; i++ {
		client.Get(ctx, "/test", nil, nil)
	}

	// Verify circuit is open
	_, err := client.Get(ctx, "/test", nil, nil)
	if err == nil {
		t.Fatal("Expected circuit breaker to be open")
	}

	// Wait for circuit timeout AND for server to recover
	time.Sleep(500 * time.Millisecond)

	// Now server is healthy and circuit should allow one request (half-open)
	resp, err := client.Get(ctx, "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected request to succeed in half-open state: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after recovery, got %d", resp.StatusCode)
	}

	// Circuit should now be closed again
	resp, err = client.Get(ctx, "/test", nil, nil)
	if err != nil {
		t.Errorf("Expected circuit to be closed after successful half-open request: %v", err)
	}
}
