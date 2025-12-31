package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestRequestInterceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Errorf("Expected X-Custom-Header to be set")
		}
		json.NewEncoder(w).Encode(TestUser{ID: 1})
	}))
	defer server.Close()

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		AddRequestInterceptor(func(req *http.Request) (*http.Request, error) {
			req.Header.Set("X-Custom-Header", "test-value")
			return req, nil
		})

	var user TestUser
	_, err := client.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestResponseInterceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(TestUser{ID: 1})
	}))
	defer server.Close()

	interceptorCalled := false

	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		AddResponseInterceptor(func(resp *http.Response) (*http.Response, error) {
			interceptorCalled = true
			return resp, nil
		})

	var user TestUser
	_, err := client.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !interceptorCalled {
		t.Error("Expected response interceptor to be called")
	}
}
