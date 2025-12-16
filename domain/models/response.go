package models

import "net/http"

// Response represents the HTTP response wrapper that GoFetch returns.
// This domain model encapsulates all response information.
type Response struct {
	StatusCode int
	Headers    http.Header
	Data       interface{}
	RawBody    []byte
}

// NewResponse creates a new Response instance.
func NewResponse(statusCode int, headers http.Header, data interface{}, rawBody []byte) *Response {
	return &Response{
		StatusCode: statusCode,
		Headers:    headers,
		Data:       data,
		RawBody:    rawBody,
	}
}
