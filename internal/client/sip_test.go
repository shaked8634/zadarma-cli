package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			"sips":[
				{"id":"123456","display_name":"SIP","lines":3},
				{"id":"649223","display_name":"SIP2","lines":2}
			],
			"left":5
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	sips, err := client.GetSIPs()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sips) != 2 {
		t.Errorf("Expected 2 SIPs, got %d", len(sips))
	}

	if sips[0]["id"] != "123456" {
		t.Errorf("Expected id 123456, got %v", sips[0]["id"])
	}

	if sips[0]["display_name"] != "SIP" {
		t.Errorf("Expected display_name SIP, got %v", sips[0]["display_name"])
	}
}

// TestGetSIPStatus tests the GetSIPStatus endpoint.
func TestGetSIPStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sip/user1/status/" {
			t.Errorf("Expected path /v1/sip/user1/status/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"sip":"user1",
			"is_online":"true"
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	isOnline, err := client.GetSIPStatus("user1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !isOnline {
		t.Errorf("Expected is_online true, got %v", isOnline)
	}
}
