package tests

import (
	"testing"
	"time"

	"github.com/fourth-ally/gofetch/infrastructure"
)

func TestClientCreation(t *testing.T) {
	client := infrastructure.NewClient()
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	if client.Config().Timeout != 30*time.Second {
		t.Errorf("Expected default timeout of 30s, got %v", client.Config().Timeout)
	}
}

func TestFluentConfiguration(t *testing.T) {
	client := infrastructure.NewClient().
		SetBaseURL("https://api.example.com").
		SetTimeout(10*time.Second).
		SetHeader("Authorization", "Bearer token")

	if client.Config().BaseURL != "https://api.example.com" {
		t.Errorf("Expected base URL to be set")
	}

	if client.Config().Timeout != 10*time.Second {
		t.Errorf("Expected timeout to be 10s")
	}

	if client.Config().Headers["Authorization"] != "Bearer token" {
		t.Errorf("Expected Authorization header to be set")
	}
}

func TestNewInstance(t *testing.T) {
	baseClient := infrastructure.NewClient().
		SetBaseURL("https://api.example.com").
		SetHeader("X-App-Version", "1.0.0").
		SetTimeout(30 * time.Second)

	derivedClient := baseClient.NewInstance().
		SetHeader("Authorization", "Bearer token")

	// Check that derived client has base settings
	if derivedClient.Config().BaseURL != "https://api.example.com" {
		t.Error("Expected derived client to inherit base URL")
	}

	if derivedClient.Config().Headers["X-App-Version"] != "1.0.0" {
		t.Error("Expected derived client to inherit X-App-Version header")
	}

	if derivedClient.Config().Timeout != 30*time.Second {
		t.Error("Expected derived client to inherit timeout")
	}

	// Check that derived client has its own settings
	if derivedClient.Config().Headers["Authorization"] != "Bearer token" {
		t.Error("Expected derived client to have Authorization header")
	}

	// Check that base client is not affected
	if _, ok := baseClient.Config().Headers["Authorization"]; ok {
		t.Error("Expected base client to not have Authorization header")
	}
}
