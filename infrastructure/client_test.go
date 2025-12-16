package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nikoskechris/gofetch/domain/errors"
)

type TestUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func TestClientCreation(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout of 30s, got %v", client.config.Timeout)
	}
}

func TestFluentConfiguration(t *testing.T) {
	client := NewClient().
		SetBaseURL("https://api.example.com").
		SetTimeout(10*time.Second).
		SetHeader("Authorization", "Bearer token")

	if client.config.BaseURL != "https://api.example.com" {
		t.Errorf("Expected base URL to be set")
	}

	if client.config.Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s")
	}

	if client.config.Headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Authorization header to be set")
	}
}

func TestGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		user := TestUser{ID: 1, Name: "John Doe", Email: "john@example.com"}
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)

	var user TestUser
	resp, err := client.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %s", user.Name)
	}
}

func TestPostRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var user TestUser
		json.NewDecoder(r.Body).Decode(&user)

		if user.Name != "Jane Doe" {
			t.Errorf("Expected name 'Jane Doe', got %s", user.Name)
		}

		user.ID = 2
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)

	newUser := TestUser{Name: "Jane Doe", Email: "jane@example.com"}
	var createdUser TestUser

	resp, err := client.Post(context.Background(), "/users", nil, newUser, &createdUser)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if createdUser.ID != 2 {
		t.Errorf("Expected ID 2, got %d", createdUser.ID)
	}
}

func TestPathParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123" {
			t.Errorf("Expected path /users/123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(TestUser{ID: 123, Name: "Test User"})
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)

	params := map[string]interface{}{
		"id": 123,
	}

	var user TestUser
	_, err := client.Get(context.Background(), "/users/:id", params, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != 123 {
		t.Errorf("Expected ID 123, got %d", user.ID)
	}
}

func TestQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("page") != "2" {
			t.Errorf("Expected page=2, got %s", query.Get("page"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", query.Get("limit"))
		}

		users := []TestUser{{ID: 1, Name: "User 1"}}
		json.NewEncoder(w).Encode(users)
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)

	params := map[string]interface{}{
		"page":  2,
		"limit": 10,
	}

	var users []TestUser
	_, err := client.Get(context.Background(), "/users", params, &users)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}
}

func TestErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)

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

func TestRequestInterceptor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Errorf("Expected X-Custom-Header to be set")
		}
		json.NewEncoder(w).Encode(TestUser{ID: 1})
	}))
	defer server.Close()

	client := NewClient().
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

	client := NewClient().
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

func TestCustomStatusValidator(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(TestUser{ID: 1})
	}))
	defer server.Close()

	// Client that only accepts 200 status
	client := NewClient().
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
	client2 := NewClient().
		SetBaseURL(server.URL).
		SetStatusValidator(func(statusCode int) bool {
			return statusCode >= 200 && statusCode < 300
		})

	_, err = client2.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestNewInstance(t *testing.T) {
	baseClient := NewClient().
		SetBaseURL("https://api.example.com").
		SetHeader("X-App-Version", "1.0.0").
		SetTimeout(30 * time.Second)

	derivedClient := baseClient.NewInstance().
		SetHeader("Authorization", "Bearer token")

	// Check that derived client has base settings
	if derivedClient.config.BaseURL != "https://api.example.com" {
		t.Error("Expected derived client to inherit base URL")
	}

	if derivedClient.config.Headers["X-App-Version"] != "1.0.0" {
		t.Error("Expected derived client to inherit X-App-Version header")
	}

	if derivedClient.config.Timeout != 30*time.Second {
		t.Error("Expected derived client to inherit timeout")
	}

	// Check that derived client has its own settings
	if derivedClient.config.Headers["Authorization"] != "Bearer token" {
		t.Error("Expected derived client to have Authorization header")
	}

	// Check that base client is not affected
	if _, ok := baseClient.config.Headers["Authorization"]; ok {
		t.Error("Expected base client to not have Authorization header")
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(TestUser{ID: 1})
	}))
	defer server.Close()

	client := NewClient().SetBaseURL(server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var user TestUser
	_, err := client.Get(ctx, "/users/1", nil, &user)
	if err == nil {
		t.Fatal("Expected context deadline exceeded error")
	}
}
