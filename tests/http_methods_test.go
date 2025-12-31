package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		user := TestUser{ID: 1, Name: "Nikos Doe", Email: "nikos@example.com"}
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)

	var user TestUser
	resp, err := client.Get(context.Background(), "/users/1", nil, &user)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if user.Name != "Nikos Doe" {
		t.Errorf("Expected name 'Nikos Doe', got %s", user.Name)
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

	client := infrastructure.NewClient().SetBaseURL(server.URL)

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

func TestPutRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		var user TestUser
		json.NewDecoder(r.Body).Decode(&user)
		user.ID = 1
		json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)
	updatedUser := TestUser{Name: "Updated Name", Email: "updated@example.com"}
	var result TestUser

	resp, err := client.Put(context.Background(), "/users/1", nil, updatedUser, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if result.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", result.Name)
	}
}

func TestPatchRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}

		var updates map[string]interface{}
		json.NewDecoder(r.Body).Decode(&updates)

		response := map[string]interface{}{"updated": true, "changes": updates}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)
	updates := map[string]string{"email": "patched@example.com"}
	var result map[string]interface{}

	resp, err := client.Patch(context.Background(), "/users/1", nil, updates, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	if result["updated"] != true {
		t.Errorf("Expected updated=true in response")
	}
}

func TestDeleteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)

	resp, err := client.Delete(context.Background(), "/users/1", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}
