package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fourth-ally/gofetch/domain/errors"
	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)

	var user TestUser
	_, err := client.Get(context.Background(), "/users/999", nil, &user)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	httpErr, ok := err.(*errors.HTTPError)
	if !ok {
		t.Fatalf("Expected HTTPError, got %T", err)
	}

	if httpErr.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", httpErr.StatusCode)
	}

	if string(httpErr.Body) != "User not found" {
		t.Errorf("Expected body 'User not found', got %s", string(httpErr.Body))
	}
}

func TestCustomStatusValidator(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	// Client that only accepts 200 status
	client := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetStatusValidator(func(statusCode int) bool {
			return statusCode == 200
		})

	var user TestUser
	_, err := client.Get(context.Background(), "/users/1", nil, &user)
	if err == nil {
		t.Fatal("Expected error for status 201, got nil")
	}

	// Client that accepts 2xx status
	client2 := infrastructure.NewClient().
		SetBaseURL(server.URL).
		SetStatusValidator(func(statusCode int) bool {
			return statusCode >= 200 && statusCode < 300
		})

	_, err = client2.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
