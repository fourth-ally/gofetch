// Package gofetch provides a robust, developer-friendly HTTP client library
// inspired by Axios, with support for interceptors, automatic JSON handling,
// and WebAssembly compatibility.
package gofetch

import (
	"github.com/nikoskechris/gofetch/infrastructure"
)

// NewClient creates a new GoFetch client instance.
// This is the primary entry point for the library.
//
// Example:
//
//	client := gofetch.NewClient().
//	    SetBaseURL("https://api.example.com").
//	    SetTimeout(10 * time.Second).
//	    SetHeader("Authorization", "Bearer token")
func NewClient() *infrastructure.Client {
	return infrastructure.NewClient()
}
