package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(TestUser{ID: 1})
	}))
	defer server.Close()

	client := infrastructure.NewClient().SetBaseURL(server.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var user TestUser
	_, err := client.Get(ctx, "/users/1", nil, &user)
	if err == nil {
		t.Fatal("Expected context deadline exceeded error")
	}
}
