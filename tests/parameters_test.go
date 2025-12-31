package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestPathParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/123" {
			t.Errorf("Expected path /users/123, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(TestUser{ID: 123, Name: "Test User"})
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)

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

	client := infrastructure.NewClient().SetBaseURL(server.URL)

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
