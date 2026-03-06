package commands

import (
	"strings"
	"testing"
)

func TestListenHasWebhookFlag(t *testing.T) {
	cmd := NewSMSCmd(nil)
	if cmd == nil {
		t.Fatal("NewSMSCmd returned nil")
	}

	// Use cobra's API: find subcommand with Use == "listen"
	found := false
	for _, c := range cmd.Commands() {
		if c.Use == "listen" {
			if c.Flags().Lookup("webhook") == nil {
				t.Fatalf("listen command missing 'webhook' flag")
			}
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("sms command does not contain 'listen' subcommand")
	}
}

func TestSetWebhookCommandExists(t *testing.T) {
	cmd := NewSMSCmd(nil)
	found := false
	for _, c := range cmd.Commands() {
		// Use HasPrefix because Use may include the argument placeholder, e.g. "set-webhook <WEBHOOK>"
		if strings.HasPrefix(c.Use, "set-webhook") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("sms command missing 'set-webhook' subcommand")
	}
}
