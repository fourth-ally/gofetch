package models

import "time"

// Config represents the configuration for the HTTP client.
// This is the domain model for client configuration.
type Config struct {
	BaseURL         string
	Timeout         time.Duration
	Headers         map[string]string
	StatusValidator func(int) bool
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

	return &Config{
		BaseURL:         c.BaseURL,
		Timeout:         c.Timeout,
		Headers:         headers,
		StatusValidator: c.StatusValidator,
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

	return merged
}
