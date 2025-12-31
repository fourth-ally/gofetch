package models

import "time"

// Config represents the configuration for the HTTP client.
// This is the domain model for client configuration.
type Config struct {
	BaseURL         string
	Timeout         time.Duration
	Headers         map[string]string
	StatusValidator func(int) bool
	RetryOptions    *RetryOptions
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		Headers:         make(map[string]string),
		Timeout:         30 * time.Second,
		StatusValidator: DefaultStatusValidator,
	}
}

// DefaultStatusValidator validates that the status code is in the 2xx range.
func DefaultStatusValidator(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// Clone creates a deep copy of the Config.
func (c *Config) Clone() *Config {
	headers := make(map[string]string, len(c.Headers))
	for k, v := range c.Headers {
		headers[k] = v
	}

	var retryOpts *RetryOptions
	if c.RetryOptions != nil {
		retryOptsCopy := *c.RetryOptions
		if len(c.RetryOptions.RetryOnStatusCodes) > 0 {
			retryOptsCopy.RetryOnStatusCodes = make([]int, len(c.RetryOptions.RetryOnStatusCodes))
			copy(retryOptsCopy.RetryOnStatusCodes, c.RetryOptions.RetryOnStatusCodes)
		}
		retryOpts = &retryOptsCopy
	}

	return &Config{
		BaseURL:         c.BaseURL,
		Timeout:         c.Timeout,
		Headers:         headers,
		StatusValidator: c.StatusValidator,
		RetryOptions:    retryOpts,
	}
}

// Merge merges another config into this one, with the other config taking precedence.
func (c *Config) Merge(other *Config) *Config {
	merged := c.Clone()

	if other.BaseURL != "" {
		merged.BaseURL = other.BaseURL
	}

	if other.Timeout != 0 {
		merged.Timeout = other.Timeout
	}

	for k, v := range other.Headers {
		merged.Headers[k] = v
	}

	if other.StatusValidator != nil {
		merged.StatusValidator = other.StatusValidator
	}

	if other.RetryOptions != nil {
		merged.RetryOptions = other.RetryOptions
	}

	return merged
}
