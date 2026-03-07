package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetDirectNumbers tests the GetDirectNumbers endpoint.
func TestGetDirectNumbers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/" {
			t.Errorf("Expected path /v1/direct_numbers/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"info":[
				{"number":"123456789012","type":"voice","status":"on","country":"Israel"},
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

	if dids[0]["number"] != "123456789012" {
		t.Errorf("Expected number 123456789012, got %v", dids[0]["number"])
	}
}

// TestGetDirectCountries tests the GetDirectCountries endpoint.
func TestGetDirectCountries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/countries/" {
			t.Errorf("Expected path /v1/direct_numbers/countries/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"info":[
				{"countryCode":"1","countryCodeIso":"US","name":"United States"},
				{"countryCode":"44","countryCodeIso":"GB","name":"United Kingdom"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	countries, err := client.GetDirectCountries()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(countries) != 2 {
		t.Errorf("Expected 2 countries, got %d", len(countries))
	}

	if countries[0]["countryCodeIso"] != "US" {
		t.Errorf("Expected countryCodeIso US, got %v", countries[0]["countryCodeIso"])
	}
}

// TestGetDirectCountry tests the GetDirectCountry endpoint.
func TestGetDirectCountry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/country/" {
			t.Errorf("Expected path /v1/direct_numbers/country/, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("country") != "US" {
			t.Errorf("Expected country US, got %s", r.URL.Query().Get("country"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"info":[
				{"id":"212","name":"New York","type":"city"}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	destinations, err := client.GetDirectCountry("US")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(destinations) != 1 {
		t.Errorf("Expected 1 destination, got %d", len(destinations))
	}

	if destinations[0]["name"] != "New York" {
		t.Errorf("Expected name 'New York', got %v", destinations[0]["name"])
	}
}

// TestGetDirectNumber tests the GetDirectNumber endpoint.
func TestGetDirectNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/direct_numbers/number/" {
			t.Errorf("Expected path /v1/direct_numbers/number/, got %s", r.URL.Path)
		}

		if r.URL.Query().Get("number") != "14155555555" {
			t.Errorf("Expected number 14155555555, got %s", r.URL.Query().Get("number"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"status":"success",
			"data":{"number":"14155555555","status":"available","price":"1.99"}
		}`))
	}))
	defer server.Close()

	client := NewClient("test_key", "test_secret", false)
	client.baseURL = server.URL + "/v1"

	info, err := client.GetDirectNumber("14155555555")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if info["number"] != "14155555555" {
		t.Errorf("Expected number 14155555555, got %v", info["number"])
	}

	if info["status"] != "available" {
		t.Errorf("Expected status available, got %v", info["status"])
	}
}
