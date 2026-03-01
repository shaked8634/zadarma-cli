package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetPBXInfo tests the GetPBXInfo endpoint.
func TestGetPBXInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pbx/internal/" {
			t.Errorf("Expected path /v1/pbx/internal/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"name":"My PBX","status":"active","pbx_id":"pbx123"}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	pbxInfo, err := client.GetPBXInfo("", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if pbxInfo["name"] != "My PBX" {
		t.Errorf("Expected name 'My PBX', got %v", pbxInfo["name"])
	}

	if pbxInfo["status"] != "active" {
		t.Errorf("Expected status active, got %v", pbxInfo["status"])
	}
}

// TestGetPBXInternalStatus tests the GetPBXInternalStatus endpoint.
func TestGetPBXInternalStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pbx/internal/pbx123/status/" {
			t.Errorf("Expected path /v1/pbx/internal/pbx123/status/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"pbx_id":"pbx123","is_online":"true","channels_active":2}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	status, err := client.GetPBXInternalStatus("pbx123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if status["pbx_id"] != "pbx123" {
		t.Errorf("Expected pbx_id pbx123, got %v", status["pbx_id"])
	}

	if status["is_online"] != "true" {
		t.Errorf("Expected is_online true, got %v", status["is_online"])
	}
}

// TestGetPBXInternalInfo tests the GetPBXInternalInfo endpoint.
func TestGetPBXInternalInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pbx/internal/pbx123/info/" {
			t.Errorf("Expected path /v1/pbx/internal/pbx123/info/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"pbx_id":"pbx123","title":"Main Office","extensions":5}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	info, err := client.GetPBXInternalInfo("pbx123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info["pbx_id"] != "pbx123" {
		t.Errorf("Expected pbx_id pbx123, got %v", info["pbx_id"])
	}

	if info["title"] != "Main Office" {
		t.Errorf("Expected title 'Main Office', got %v", info["title"])
	}
}

// TestSetWebhook tests the SetWebhook endpoint.
func TestSetWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pbx/webhooks/url/" {
			t.Errorf("Expected path /v1/pbx/webhooks/url/, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check Content-Type header
		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type application/x-www-form-urlencoded, got %s", ct)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"url":"https://example.com/webhook","status":"set"}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	result, err := client.SetWebhook("https://example.com/webhook")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["status"] != "set" {
		t.Errorf("Expected status 'set', got %v", result["status"])
	}
}

// TestGetWebhook tests the GetWebhooks endpoint.
func TestGetWebhook(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pbx/webhooks/url/" {
			t.Errorf("Expected path /v1/pbx/webhooks/url/, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"url":"https://example.com/webhook","status":"active"}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	webhook, err := client.GetWebhooks()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if webhook["url"] != "https://example.com/webhook" {
		t.Errorf("Expected url 'https://example.com/webhook', got %v", webhook["url"])
	}

	if webhook["status"] != "active" {
		t.Errorf("Expected status 'active', got %v", webhook["status"])
	}
}
