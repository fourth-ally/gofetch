package infrastructure

import (
	"math"
	"math/rand"
	"time"

	"github.com/fourth-ally/gofetch/domain/models"
)

// RetryManager handles retry logic with backoff strategies.
type RetryManager struct {
	options *models.RetryOptions
	rng     *rand.Rand
}

// NewRetryManager creates a new retry manager.
func NewRetryManager(options *models.RetryOptions) *RetryManager {
	if options == nil {
		options = models.NewRetryOptions()
	}

	return &RetryManager{
		options: options,
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ShouldRetry determines if a request should be retried based on the error and attempt number.
func (rm *RetryManager) ShouldRetry(attempt int, statusCode int, err error) bool {
	// Don't retry if max retries reached
	if attempt >= rm.options.MaxRetries {
		return false
	}

	// Retry on network errors
	if err != nil {
		return true
	}

	// Retry on configured status codes
	return rm.options.ShouldRetryStatus(statusCode)
}

// CalculateDelay calculates the delay before the next retry attempt.
func (rm *RetryManager) CalculateDelay(attempt int) time.Duration {
	var delay time.Duration

	switch rm.options.Backoff {
	case models.BackoffExponential:
		delay = rm.calculateExponentialBackoff(attempt)
	case models.BackoffLinear:
		delay = rm.calculateLinearBackoff(attempt)
	case models.BackoffFixed:
		delay = rm.options.InitialDelay
	default:
		delay = rm.calculateExponentialBackoff(attempt)
	}

	// Apply max delay cap
	if delay > rm.options.MaxDelay {
		delay = rm.options.MaxDelay
	}

	// Apply jitter if enabled
	if rm.options.Jitter {
		delay = rm.applyJitter(delay)
	}

	return delay
}

// calculateExponentialBackoff calculates exponential backoff: initialDelay * 2^attempt.
func (rm *RetryManager) calculateExponentialBackoff(attempt int) time.Duration {
	multiplier := math.Pow(2, float64(attempt))
	delay := time.Duration(float64(rm.options.InitialDelay) * multiplier)
	return delay
}

// calculateLinearBackoff calculates linear backoff: initialDelay * (attempt + 1).
func (rm *RetryManager) calculateLinearBackoff(attempt int) time.Duration {
	return rm.options.InitialDelay * time.Duration(attempt+1)
}

// applyJitter adds random jitter to the delay.
func (rm *RetryManager) applyJitter(delay time.Duration) time.Duration {
	jitterFraction := rm.options.JitterFraction
	if jitterFraction <= 0 {
		jitterFraction = 0.3 // Default 30%
	}
	if jitterFraction > 1 {
		jitterFraction = 1.0
	}

	// Calculate jitter range: delay ± (delay * jitterFraction)
	jitterRange := float64(delay) * jitterFraction
	randomJitter := rm.rng.Float64()*2*jitterRange - jitterRange

	newDelay := float64(delay) + randomJitter
	if newDelay < 0 {
		newDelay = 0
	}

	return time.Duration(newDelay)
}

// Wait pauses execution for the calculated delay.
func (rm *RetryManager) Wait(attempt int) {
	delay := rm.CalculateDelay(attempt)
	time.Sleep(delay)
}
