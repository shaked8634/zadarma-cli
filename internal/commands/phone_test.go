package commands

import (
	"testing"
)

// TestPhoneCommand verifies phone command structure.
func TestPhoneCommand(t *testing.T) {
	// Test that NewPhoneCmd returns a command
	// Using nil factory since we only check command structure
	cmd := NewPhoneCmd(nil)
	if cmd == nil {
		t.Error("Expected NewPhoneCmd to return a command, got nil")
	}

	if cmd.Use != "phone" {
		t.Errorf("Expected command name 'phone', got %s", cmd.Use)
	}

	// Check that subcommands exist
	subcommands := []string{"list", "countries", "country", "number"}
	for _, subCmd := range subcommands {
		_, _, err := cmd.Find([]string{subCmd})
		if err != nil {
			t.Errorf("Expected '%s' subcommand to exist, got error: %v", subCmd, err)
		}
	}
}
