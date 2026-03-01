package commands

import (
	"testing"
)

// TestSIPCommand verifies sip command structure.
func TestSIPCommand(t *testing.T) {
	// Test that NewSIPCmd returns a command
	// Using nil factory since we only check command structure
	cmd := NewSIPCmd(nil)
	if cmd == nil {
		t.Error("Expected NewSIPCmd to return a command, got nil")
	}

	if cmd.Use != "sip" {
		t.Errorf("Expected command name 'sip', got %s", cmd.Use)
	}
}
