package errors

import (
	"fmt"
	"net/http"
)

// HTTPError represents an error response from an HTTP request.
// This is the domain model for all HTTP-related errors in the system.
type HTTPError struct {
	StatusCode   int
	Body         []byte
	Headers      http.Header
	Message      string
	OriginalResp *http.Response
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

// NewHTTPError creates a new HTTPError from an HTTP response.
func NewHTTPError(resp *http.Response, body []byte, message string) *HTTPError {
	return &HTTPError{
		StatusCode:   resp.StatusCode,
		Body:         body,
		Headers:      resp.Header,
		Message:      message,
		OriginalResp: resp,
	}
}
