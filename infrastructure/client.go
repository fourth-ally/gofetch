package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nkechr1s/gofetch/domain/contracts"
	"github.com/nkechr1s/gofetch/domain/errors"
	"github.com/nkechr1s/gofetch/domain/models"
)

// Client is the main HTTP client implementation.
type Client struct {
	httpClient           *http.Client
	config               *models.Config
	requestInterceptors  []contracts.RequestInterceptor
	responseInterceptors []contracts.ResponseInterceptor
	dataTransformer      contracts.DataTransformer
	uploadProgress       contracts.ProgressCallback
	downloadProgress     contracts.ProgressCallback
}

// NewClient creates a new GoFetch client instance.
func NewClient() *Client {
	return &Client{
		httpClient:           &http.Client{Timeout: 30 * time.Second},
		config:               models.NewConfig(),
		requestInterceptors:  make([]contracts.RequestInterceptor, 0),
		responseInterceptors: make([]contracts.ResponseInterceptor, 0),
	}
}

// SetBaseURL sets the base URL for all requests.
func (c *Client) SetBaseURL(baseURL string) *Client {
	c.config.BaseURL = baseURL
	return c
}

// SetTimeout sets the timeout for requests.
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.config.Timeout = timeout
	c.httpClient.Timeout = timeout
	return c
}

// SetHeader sets a default header for all requests.
func (c *Client) SetHeader(key, value string) *Client {
	c.config.Headers[key] = value
	return c
}

// SetStatusValidator sets a custom status validator function.
func (c *Client) SetStatusValidator(validator func(int) bool) *Client {
	c.config.StatusValidator = validator
	return c
}

// AddRequestInterceptor adds a request interceptor.
func (c *Client) AddRequestInterceptor(interceptor contracts.RequestInterceptor) *Client {
	c.requestInterceptors = append(c.requestInterceptors, interceptor)
	return c
}

// AddResponseInterceptor adds a response interceptor.
func (c *Client) AddResponseInterceptor(interceptor contracts.ResponseInterceptor) *Client {
	c.responseInterceptors = append(c.responseInterceptors, interceptor)
	return c
}

// SetDataTransformer sets the data transformer function.
func (c *Client) SetDataTransformer(transformer contracts.DataTransformer) *Client {
	c.dataTransformer = transformer
	return c
}

// SetUploadProgress sets the upload progress callback.
func (c *Client) SetUploadProgress(callback contracts.ProgressCallback) *Client {
	c.uploadProgress = callback
	return c
}

// SetDownloadProgress sets the download progress callback.
func (c *Client) SetDownloadProgress(callback contracts.ProgressCallback) *Client {
	c.downloadProgress = callback
	return c
}

// NewInstance creates a new client instance inheriting all settings from the current client.
func (c *Client) NewInstance() *Client {
	newClient := &Client{
		httpClient:           &http.Client{Timeout: c.config.Timeout},
		config:               c.config.Clone(),
		requestInterceptors:  make([]contracts.RequestInterceptor, len(c.requestInterceptors)),
		responseInterceptors: make([]contracts.ResponseInterceptor, len(c.responseInterceptors)),
		dataTransformer:      c.dataTransformer,
		uploadProgress:       c.uploadProgress,
		downloadProgress:     c.downloadProgress,
	}

	copy(newClient.requestInterceptors, c.requestInterceptors)
	copy(newClient.responseInterceptors, c.responseInterceptors)

	return newClient
}

// buildURL constructs the full URL from base URL, path, and parameters.
func (c *Client) buildURL(path string, params map[string]interface{}) (string, error) {
	// Start with base URL or empty string
	fullURL := c.config.BaseURL

	// Handle path parameters (e.g., /users/:id)
	processedPath := path
	queryParams := url.Values{}

	if params != nil {
		for key, value := range params {
			placeholder := ":" + key
			if strings.Contains(processedPath, placeholder) {
				// Replace path parameter
				processedPath = strings.Replace(processedPath, placeholder, fmt.Sprintf("%v", value), -1)
			} else {
				// Add to query string
				queryParams.Add(key, fmt.Sprintf("%v", value))
			}
		}
	}

	// Combine base URL and path
	if fullURL != "" && processedPath != "" {
		fullURL = strings.TrimRight(fullURL, "/") + "/" + strings.TrimLeft(processedPath, "/")
	} else if processedPath != "" {
		fullURL = processedPath
	}

	// Add query parameters
	if len(queryParams) > 0 {
		fullURL += "?" + queryParams.Encode()
	}

	return fullURL, nil
}

// executeRequest executes an HTTP request with all interceptors and error handling.
func (c *Client) executeRequest(ctx context.Context, method, path string, params map[string]interface{}, body interface{}, target interface{}, requestConfig *models.Config) (*models.Response, error) {
	// Merge configurations
	config := c.config
	if requestConfig != nil {
		config = c.config.Merge(requestConfig)
	}

	// Build URL
	fullURL, err := c.buildURL(path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	// Prepare request body
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonData)

		// Wrap with progress tracking if callback is set
		if c.uploadProgress != nil {
			bodyReader = &progressReader{
				reader:   bodyReader,
				total:    int64(len(jsonData)),
				callback: c.uploadProgress,
			}
		}
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Set content type for body requests
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply request interceptors
	for _, interceptor := range c.requestInterceptors {
		req, err = interceptor(req)
		if err != nil {
			return nil, fmt.Errorf("request interceptor error: %w", err)
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution error: %w", err)
	}
	defer resp.Body.Close()

	// Apply response interceptors
	for _, interceptor := range c.responseInterceptors {
		resp, err = interceptor(resp)
		if err != nil {
			return nil, fmt.Errorf("response interceptor error: %w", err)
		}
	}

	// Read response body with progress tracking
	var respBody []byte
	if c.downloadProgress != nil && resp.ContentLength > 0 {
		progressReader := &progressReader{
			reader:   resp.Body,
			total:    resp.ContentLength,
			callback: c.downloadProgress,
		}
		respBody, err = io.ReadAll(progressReader)
	} else {
		respBody, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Validate status code
	if !config.StatusValidator(resp.StatusCode) {
		return nil, errors.NewHTTPError(resp, respBody, "")
	}

	// Apply data transformer if set
	if c.dataTransformer != nil {
		respBody, err = c.dataTransformer(respBody)
		if err != nil {
			return nil, fmt.Errorf("data transformer error: %w", err)
		}
	}

	// Unmarshal response into target if provided
	if target != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, target); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return models.NewResponse(resp.StatusCode, resp.Header, target, respBody), nil
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, params map[string]interface{}, target interface{}) (*models.Response, error) {
	return c.executeRequest(ctx, http.MethodGet, path, params, nil, target, nil)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, params map[string]interface{}, body interface{}, target interface{}) (*models.Response, error) {
	return c.executeRequest(ctx, http.MethodPost, path, params, body, target, nil)
}

// Put performs a PUT request.
func (c *Client) Put(ctx context.Context, path string, params map[string]interface{}, body interface{}, target interface{}) (*models.Response, error) {
	return c.executeRequest(ctx, http.MethodPut, path, params, body, target, nil)
}

// Patch performs a PATCH request.
func (c *Client) Patch(ctx context.Context, path string, params map[string]interface{}, body interface{}, target interface{}) (*models.Response, error) {
	return c.executeRequest(ctx, http.MethodPatch, path, params, body, target, nil)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, params map[string]interface{}, target interface{}) (*models.Response, error) {
	return c.executeRequest(ctx, http.MethodDelete, path, params, nil, target, nil)
}
