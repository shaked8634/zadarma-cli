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

	client := NewClient("test_key", "test_secret", false)
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
// (Moved to sip_test.go)

// TestGetSIPStatus tests the GetSIPStatus endpoint.
// (Moved to sip_test.go)

// TestGetDIDs tests the GetDirectNumbers endpoint for fetching all numbers.
// (Moved to phone_test.go)

// TestGetDIDsByNumber tests the GetDirectNumbers endpoint for fetching specific numbers.
// (Moved to phone_test.go)

// TestGetDIDsByMultipleNumbers tests the GetDirectNumbers endpoint for fetching multiple specific numbers.
// (Moved to phone_test.go)
// TestSendSMS tests the SendSMS endpoint.
// (Moved to sms_test.go)

// TestGetSMSSenders tests the GetSMSSenders endpoint.
// (Moved to sms_test.go)

// TestGetPrice tests the GetPrice endpoint.
func TestGetPrice(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/info/price/" {
			t.Errorf("Expected path /v1/info/price/, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("number") != "+14155555555" {
			t.Errorf("Expected number +14155555555, got %s", r.URL.Query().Get("number"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
            "status":"success",
            "data":{"price":"0.050","currency":"USD"}
        }`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	data, err := client.GetPrice("+14155555555")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data["price"] != "0.050" {
		t.Errorf("Expected price 0.050, got %v", data["price"])
	}

	if data["currency"] != "USD" {
		t.Errorf("Expected currency USD, got %v", data["currency"])
	}
}

// TestSendSMS tests the SendSMS endpoint.
// (Moved to sms_test.go)

// TestGetSMSSenders tests the GetSMSSenders endpoint.
// (Moved to sms_test.go)

// TestGetPBXInfo tests the GetPBXInfo endpoint.
// (Moved to pbx_test.go)
