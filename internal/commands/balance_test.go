package commands

import (
	"testing"
)

// TestBalanceCommand verifies balance command structure.
func TestBalanceCommand(t *testing.T) {
	// Test that NewBalanceCmd returns a command
	// Using nil factory since we only check command structure
	cmd := NewBalanceCmd(nil)
	if cmd == nil {
		t.Fatal("Expected NewBalanceCmd to return a command, got nil")
	}

	if cmd.Use != "balance" {
		t.Errorf("Expected command name 'balance', got %s", cmd.Use)
	}
}
