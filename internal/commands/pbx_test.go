package commands

import (
	"testing"
)

// TestPBXInfoCommand verifies pbx command structure.
func TestPBXInfoCommand(t *testing.T) {
	// Test that NewPBXCmd returns a command
	// Using nil factory since we only check command structure
	cmd := NewPBXCmd(nil)
	if cmd == nil {
		t.Error("Expected NewPBXCmd to return a command, got nil")
	}

	if cmd.Use != "pbx" {
		t.Errorf("Expected command name 'pbx', got %s", cmd.Use)
	}

	// Check that 'info' subcommand exists
	infoCmd, _, err := cmd.Find([]string{"info"})
	if err != nil {
		t.Errorf("Expected 'info' subcommand to exist, got error: %v", err)
	}

	if infoCmd != nil && infoCmd.Use != "info" {
		t.Errorf("Expected subcommand name 'info', got %s", infoCmd.Use)
	}
}
