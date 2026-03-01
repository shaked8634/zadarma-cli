package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/zadarma/zadarma-cli/internal/client"
)

func TestStatisticsValidation(t *testing.T) {
	mockFactory := func() *client.Client { return client.NewClient("test", "test", false) }

	tests := []struct {
		name        string
		args        []string
		expectedErr string
	}{
		{
			name:        "invalid start date",
			args:        []string{"statistics", "--start", "3434"},
			expectedErr: "invalid start date format",
		},
		{
			name:        "invalid end date",
			args:        []string{"statistics", "--end", "2026-02-28"},
			expectedErr: "invalid end date format",
		},
		{
			name:        "invalid sip negative",
			args:        []string{"statistics", "--sip", "-1"},
			expectedErr: "invalid sip",
		},
		{
			name:        "valid dates",
			args:        []string{"statistics", "--start", "2026-02-28 10:00:00", "--end", "2026-02-28 11:00:00"},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewStatisticsCmd(mockFactory)
			cmd.SetArgs(tt.args[1:])

			// Suppress output for tests
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			// For the success case, we only want to validate flags, not perform a real API call.
			if tt.expectedErr == "" {
				cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			}

			err := cmd.Execute()
			if tt.expectedErr == "" {
				// Success if no validation error is returned
				if err != nil && (strings.Contains(err.Error(), "invalid start date format") ||
					strings.Contains(err.Error(), "invalid end date format") ||
					strings.Contains(err.Error(), "invalid sip")) {
					t.Errorf("expected validation to pass, but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, but got none", tt.expectedErr)
				} else if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Errorf("expected error containing %q, but got: %v", tt.expectedErr, err)
				}
			}
		})
	}
}
