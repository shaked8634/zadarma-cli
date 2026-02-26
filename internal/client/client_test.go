package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetBalance tests the GetBalance endpoint.
func TestGetBalance(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/info/balance/" {
			t.Errorf("Expected path /v1/info/balance/, got %s", r.URL.Path)
		}

		// Check authorization header exists
		if r.Header.Get("Authorization") == "" {
			t.Error("Authorization header is missing")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"success","balance":123.45,"currency":"USD"}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret")
	client.baseURL = server.URL + "/v1"

	balance, currency, err := client.GetBalance()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if balance != 123.45 {
		t.Errorf("Expected balance 123.45, got %v", balance)
	}

	if currency != "USD" {
		t.Errorf("Expected currency USD, got %s", currency)
	}
}

// TestGetSIPs tests the GetSIPs endpoint.
func TestGetSIPs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sip/" {
			t.Errorf("Expected path /v1/sip/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":[
				{"sip_user":"user1","status":"active","password":"pass1"},
				{"sip_user":"user2","status":"inactive","password":"pass2"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret")
	client.baseURL = server.URL + "/v1"

	sips, err := client.GetSIPs()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sips) != 2 {
		t.Errorf("Expected 2 SIPs, got %d", len(sips))
	}

	if sips[0]["sip_user"] != "user1" {
		t.Errorf("Expected sip_user user1, got %v", sips[0]["sip_user"])
	}
}

// TestGetDIDs tests the GetDIDs endpoint.
func TestGetDIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/did/" {
			t.Errorf("Expected path /v1/did/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":[
				{"number":"+14155555555","type":"voice","status":"active"},
				{"number":"+14155555556","type":"fax","status":"active"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret")
	client.baseURL = server.URL + "/v1"

	dids, err := client.GetDIDs()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dids) != 2 {
		t.Errorf("Expected 2 DIDs, got %d", len(dids))
	}

	if dids[0]["number"] != "+14155555555" {
		t.Errorf("Expected number +14155555555, got %v", dids[0]["number"])
	}
}

// TestSendSMS tests the SendSMS endpoint.
func TestSendSMS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sms/" {
			t.Errorf("Expected path /v1/sms/, got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check query parameters
		if r.URL.Query().Get("number") != "+14155555555" {
			t.Errorf("Expected number +14155555555, got %s", r.URL.Query().Get("number"))
		}

		if r.URL.Query().Get("message") != "Hello World" {
			t.Errorf("Expected message 'Hello World', got %s", r.URL.Query().Get("message"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"id":"msg123","status":"sent","timestamp":"2025-02-26T07:42:00Z"}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret")
	client.baseURL = server.URL + "/v1"

	result, err := client.SendSMS("+14155555555", "Hello World")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "msg123" {
		t.Errorf("Expected id msg123, got %v", result["id"])
	}

	if result["status"] != "sent" {
		t.Errorf("Expected status sent, got %v", result["status"])
	}
}

// TestGetPBXInfo tests the GetPBXInfo endpoint.
func TestGetPBXInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pbx/" {
			t.Errorf("Expected path /v1/pbx/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"name":"My PBX","status":"active","pbx_id":"pbx123"}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret")
	client.baseURL = server.URL + "/v1"

	pbxInfo, err := client.GetPBXInfo()
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
