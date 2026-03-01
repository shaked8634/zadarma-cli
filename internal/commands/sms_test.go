package commands

import (
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
