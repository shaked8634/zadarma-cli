package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSendSMS tests the SendSMS endpoint.
func TestSendSMS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sms/send/" {
			t.Errorf("Expected path /v1/sms/send/, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

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
		// Return the full response (what the real API returns)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"messages":1,
			"cost":0.2,
			"currency":"EUR",
			"sms_detalization":[{"callerid":"+14155555555","number":"14155555555","cost":0.2,"parts":1,"message":"Hello World"}]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false, false)
	client.baseURL = server.URL + "/v1"

	result, err := client.SendSMS("+14155555555", "Hello World", "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify full response is returned
	if result["status"] != "success" {
		t.Errorf("Expected status success, got %v", result["status"])
	}

	if messages, ok := result["messages"].(float64); !ok || messages != 1 {
		t.Errorf("Expected messages 1, got %v", result["messages"])
	}

	if cost, ok := result["cost"].(float64); !ok || cost != 0.2 {
		t.Errorf("Expected cost 0.2, got %v", result["cost"])
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

	client := NewClient("test_key", "test_secret", false, false)
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
