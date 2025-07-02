package esimgo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient("test-api-key")

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.apiKey)
	}

	if client.baseURL != "https://api.esim-go.com/v2.4" {
		t.Errorf("Expected base URL 'https://api.esim-go.com/v2.4', got '%s'", client.baseURL)
	}
}

func TestESIMGoClient(t *testing.T) {
	client := NewESIMGoClient("test-api-key")

	if client.ESIMs == nil {
		t.Error("ESIMs service should not be nil")
	}
	if client.Organization == nil {
		t.Error("Organization service should not be nil")
	}
	if client.Orders == nil {
		t.Error("Orders service should not be nil")
	}
	if client.Inventory == nil {
		t.Error("Inventory service should not be nil")
	}
	if client.Catalogue == nil {
		t.Error("Catalogue service should not be nil")
	}
	if client.Networks == nil {
		t.Error("Networks service should not be nil")
	}
}

func TestMakeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test-api-key" {
			t.Errorf("Expected X-API-Key header 'test-api-key', got '%s'", r.Header.Get("X-API-Key"))
		}

		response := map[string]string{"status": "success"}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.SetBaseURL(server.URL)

	var result map[string]string
	err := client.makeRequest(context.Background(), "GET", "/test", nil, &result)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", result["status"])
	}
}

func TestOrganizationService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/organisation" && r.Method == "GET" {
			response := Organization{
				Name:     "Test Organization",
				Currency: "USD",
				Balance:  1000,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	client := NewESIMGoClient("test-api-key")
	client.SetBaseURL(server.URL)

	org, err := client.Organization.GetDetails(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if org.Name != "Test Organization" {
		t.Errorf("Expected name 'Test Organization', got '%s'", org.Name)
	}

	if org.Balance != 1000 {
		t.Errorf("Expected balance 1000, got %d", org.Balance)
	}
}
