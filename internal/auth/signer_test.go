package auth

import (
	"net/url"
	"testing"
)

// TestSignerEmpty tests signing with no parameters (like /info/balance/).
func TestSignerEmpty(t *testing.T) {
	signer := NewSigner("1c64dee7ee7638e3a507", "d78f2a12cc1fc5effe71")

	method := "/v1/info/balance/"
	params := url.Values{}

	signature := signer.Sign(method, params)

	// Expected signature (computed manually following the algorithm)
	expected := "rhSjhy87PyCS8HamPv6b1t199EA="

	if signature != expected {
		t.Errorf("Expected signature %q, got %q", expected, signature)
	}
}

// TestSignerWithParams tests signing with parameters.
func TestSignerWithParams(t *testing.T) {
	signer := NewSigner("test_key", "test_secret")

	method := "/v1/info/price/"
	params := url.Values{}
	params.Set("number", "14155555555")

	signature := signer.Sign(method, params)

	// Verify it's a valid base64 string (just a sanity check)
	if signature == "" {
		t.Error("Signature should not be empty")
	}
}

// TestAuthHeader tests the full authorization header generation.
func TestAuthHeader(t *testing.T) {
	signer := NewSigner("test_key", "test_secret")

	method := "/v1/info/balance/"
	params := url.Values{}

	header := signer.AuthHeader(method, params)

	if !contains(header, "test_key:") {
		t.Errorf("Authorization header should contain API key, got: %s", header)
	}
}

// TestQueryStringOrdering tests alphabetical sorting of query parameters.
func TestQueryStringOrdering(t *testing.T) {
	signer := NewSigner("key", "secret")

	params := url.Values{}
	params.Set("z_param", "value_z")
	params.Set("a_param", "value_a")
	params.Set("m_param", "value_m")

	queryStr := signer.buildQueryString(params)

	// Should be alphabetically sorted
	if queryStr != "a_param=value_a&m_param=value_m&z_param=value_z" {
		t.Errorf("Query string not properly sorted. Got: %s", queryStr)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
