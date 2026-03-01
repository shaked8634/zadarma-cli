package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// wantsJSON reports whether the current command invocation requested JSON output
// via --output=json flag.
func wantsJSON(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if root := cmd.Root(); root != nil {
		if of, err := root.PersistentFlags().GetString("output"); err == nil && of == "json" {
			return true
		}
	}
	return false
}

// failCmd prints an error in the appropriate format and exits with status 1.
// When JSON output is requested, it prints: {"status":"error","message":"..."}
// Otherwise, it prints a human-friendly text error.
func failCmd(cmd *cobra.Command, err error) error { // return type only to satisfy RunE signature; never returns
	if err == nil {
		os.Exit(1)
	}
	if wantsJSON(cmd) {
		payload := map[string]any{
			"status":  "error",
			"message": err.Error(),
		}
		b, _ := json.Marshal(payload)
		// Errors are part of program output contract in JSON mode; print to stdout
		fmt.Println(string(b))
	} else {
		// In text mode keep errors on stderr to play well with pipes
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(1)
	return err // unreachable
}

// getStringField safely extracts string fields from response maps
func getStringField(m map[string]interface{}, key, defaultVal string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
		return fmt.Sprintf("%v", val)
	}
	return defaultVal
}
