package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fourth-ally/gofetch/domain/errors"
	"github.com/fourth-ally/gofetch/domain/models"
	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestDataTransformer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer server.Close()

	transformerCalled := false
	transformer := func(data []byte) ([]byte, error) {
		transformerCalled = true
		return data, nil
	}

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetDataTransformer(transformer)

	_, err := client.Get(context.Background(), "/test", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !transformerCalled {
		t.Error("Expected data transformer to be called")
	}
}

func TestUploadProgress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	progressCalled := false
	progressCallback := func(loaded, total int64) {
		progressCalled = true
	}

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetUploadProgress(progressCallback)

	body := strings.Repeat("data", 1000)
	_, err := client.Post(context.Background(), "/upload", nil, body, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !progressCalled {
		t.Error("Expected upload progress callback to be called")
	}
}

func TestDownloadProgress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "4000")
		w.Write([]byte(strings.Repeat("data", 1000)))
	}))
	defer server.Close()

	progressCalled := false
	progressCallback := func(loaded, total int64) {
		progressCalled = true
	}

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetDownloadProgress(progressCallback)

	_, err := client.Get(context.Background(), "/download", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !progressCalled {
		t.Error("Expected download progress callback to be called")
	}
}

func TestConfigMerge(t *testing.T) {
	config1 := models.NewConfig()
	config1.BaseURL = "https://api.example.com"
	config1.Timeout = 5 * time.Second
	config1.Headers = map[string]string{"Authorization": "Bearer token1"}

	config2 := &models.Config{
		BaseURL: "https://api2.example.com",
		Headers: map[string]string{"X-Custom": "value"},
	}

	merged := config1.Merge(config2)

	if merged.BaseURL != "https://api2.example.com" {
		t.Errorf("Expected BaseURL to be merged, got %s", merged.BaseURL)
	}

	if merged.Timeout != 5*time.Second {
		t.Errorf("Expected Timeout to remain 5s, got %v", merged.Timeout)
	}

	if merged.Headers["Authorization"] != "Bearer token1" {
		t.Error("Expected original header to be preserved")
	}

	if merged.Headers["X-Custom"] != "value" {
		t.Error("Expected merged header to be present")
	}
}

func TestHTTPErrorMethod(t *testing.T) {
	// Create a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "resource not found"}`))
	}))
	defer server.Close()

	// Make a request that will fail
	resp, _ := http.Get(server.URL)
	body := []byte(`{"error": "resource not found"}`)

	httpErr := errors.NewHTTPError(resp, body, "Not Found")

	errorMsg := httpErr.Error()

	if !strings.Contains(errorMsg, "404") {
		t.Errorf("Expected error message to contain status code, got: %s", errorMsg)
	}

	if !strings.Contains(errorMsg, "Not Found") {
		t.Errorf("Expected error message to contain status text, got: %s", errorMsg)
	}
}
