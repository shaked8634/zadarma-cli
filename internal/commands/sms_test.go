package commands

import (
	"io"
	"net/http"
	"testing"
)

// TestSMSCommand verifies sms command structure.
func TestSMSCommand(t *testing.T) {
	// Test that NewSMSCmd returns a command
	// Using nil factory since we only check command structure
	cmd := NewSMSCmd(nil)
	if cmd == nil {
		t.Error("Expected NewSMSCmd to return a command, got nil")
	}

	if cmd.Use != "sms" {
		t.Errorf("Expected command name 'sms', got %s", cmd.Use)
	}
}

// TestStartSMSListener verifies the behavior of startSMSListener.
func TestStartSMSListener(t *testing.T) {
	// Mock port and output mode
	port := "8080"
	jsonOutput := false

	// Run the listener in a goroutine
	go func() {
		err := startSMSListener(port, jsonOutput)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}()

	// Simulate a verification request
	resp, err := http.Get("http://localhost:" + port + "?zd_echo=test")
	if err != nil {
		t.Fatalf("Failed to send verification request: %v", err)
	}
	defer resp.Body.Close()

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
