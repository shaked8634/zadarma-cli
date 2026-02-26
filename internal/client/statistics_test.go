package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/statistics/" {
			t.Errorf("Expected path /v1/statistics/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":[
				{"callstart":"2025-02-26 10:00:00","sip":"12345","billseconds":60,"billcost":0.05,"destination":"+123456789"},
				{"callstart":"2025-02-26 11:00:00","sip":"12345","billseconds":120,"billcost":0.10,"destination":"+987654321"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret")
	client.baseURL = server.URL + "/v1"

	stats, err := client.GetStatistics(nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(stats) != 2 {
		t.Errorf("Expected 2 stats, got %d", len(stats))
	}

	if stats[0]["billseconds"].(float64) != 60 {
		t.Errorf("Expected 60 seconds, got %v", stats[0]["billseconds"])
	}
}
