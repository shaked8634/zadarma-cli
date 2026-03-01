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

	client := NewClient("test_key", "test_secret", false)
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

// TestGetDIDs tests the GetDirectNumbers endpoint for fetching all numbers.
func TestGetDIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/" {
			t.Errorf("Expected path /v1/direct_numbers/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"info":[
				{"number":"972556620707","type":"voice","status":"on","country":"Israel"},
				{"number":"19293091254","type":"fax","status":"on","country":"United States"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	dids, err := client.GetDirectNumbers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dids) != 2 {
		t.Errorf("Expected 2 DIDs, got %d", len(dids))
	}

	if dids[0]["number"] != "972556620707" {
		t.Errorf("Expected number 972556620707, got %v", dids[0]["number"])
	}
}

// TestGetDIDsByNumber tests the GetDirectNumbers endpoint for fetching specific numbers.
func TestGetDIDsByNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/" {
			t.Errorf("Expected path /v1/direct_numbers/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"info":[
				{"number":"972556620707","type":"voice","status":"on","country":"Israel"},
				{"number":"19293091254","type":"fax","status":"on","country":"United States"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	dids, err := client.GetDirectNumbers("972556620707")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dids) != 1 {
		t.Errorf("Expected 1 DID, got %d", len(dids))
	}

	if dids[0]["number"] != "972556620707" {
		t.Errorf("Expected number 972556620707, got %v", dids[0]["number"])
	}
}

// TestGetDIDsByMultipleNumbers tests the GetDirectNumbers endpoint for fetching multiple specific numbers.
func TestGetDIDsByMultipleNumbers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/" {
			t.Errorf("Expected path /v1/direct_numbers/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"info":[
				{"number":"972556620707","type":"voice","status":"on","country":"Israel"},
				{"number":"19293091254","type":"fax","status":"on","country":"United States"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	dids, err := client.GetDirectNumbers("972556620707", "19293091254")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(dids) != 2 {
		t.Errorf("Expected 2 DIDs, got %d", len(dids))
	}
}

// TestSendSMS tests the SendSMS endpoint.
func TestSendSMS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sms/send/" {
			t.Errorf("Expected path /v1/sms/send/, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Parse form body and check parameters
		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		if r.PostForm.Get("number") != "+14155555555" {
			t.Errorf("Expected number +14155555555, got %s", r.PostForm.Get("number"))
		}

		if r.PostForm.Get("message") != "Hello World" {
			t.Errorf("Expected message 'Hello World', got %s", r.PostForm.Get("message"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
            "status":"success",
            "data":{"id":"msg123","status":"sent","timestamp":"2025-02-26T07:42:00Z"}
        }`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	result, err := client.SendSMS("+14155555555", "Hello World", "")
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

// TestGetSMSSenders tests the GetSMSSenders endpoint.
func TestGetSMSSenders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sms/senderid/" {
			t.Errorf("Expected path /v1/sms/senderid/, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("phones") != "+14155555555" {
			t.Errorf("Expected phones +14155555555, got %s", r.URL.Query().Get("phones"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":[
				{"sender_id":"Zadarma","type":"alpha"},
				{"sender_id":"14155551111","type":"number"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	senders, err := client.GetSMSSenders("+14155555555")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(senders) != 2 {
		t.Errorf("Expected 2 senders, got %d", len(senders))
	}

	if senders[0]["sender_id"] != "Zadarma" {
		t.Errorf("Expected sender_id Zadarma, got %v", senders[0]["sender_id"])
	}
}

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
