package commands

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestStartSMSListener verifies the behavior of startSMSListener.
func TestStartSMSListener(t *testing.T) {
	// Use a random available port
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("Failed to close listener: %v", err)
	}

	portStr := fmt.Sprintf("%d", port)

	jsonOutput := false

	// Run the listener in a goroutine
	go func() {
		err := startSMSListener(portStr, jsonOutput, nil)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Simulate a verification request
	resp, err := http.Get("http://localhost:" + portStr + "?zd_echo=test")
	if err != nil {
		t.Fatalf("Failed to send verification request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "test" {
		t.Errorf("Expected response 'test', got %s", string(body))
	}
}
