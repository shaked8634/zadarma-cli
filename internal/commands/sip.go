package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/log"
)

// NewSIPCmd creates the 'sip' command group.
func NewSIPCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sip",
		Short: "Manage SIP accounts",
		Long:  "SIP account management commands",
		// Do not show usage on API/runtime errors
		SilenceUsage: true,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:          "list",
			Short:        "List SIP accounts",
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput := wantsJSON(cmd)
				c := factory()

				sips, err := c.GetSIPs()
				if err != nil {
					log.Debugf("Failed to get SIP accounts: %v", err)
					return failCmd(cmd, err)
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(sips, "", "  ")
					fmt.Println(string(out))
				} else {
					if len(sips) == 0 {
						fmt.Println("No SIP accounts found.")
						return nil
					}
					printSIPTable(sips)
				}
				return nil
			},
		},
		&cobra.Command{
			Use:          "info <ID>",
			Short:        "Get SIP account info",
			Args:         cobra.ExactArgs(1),
			SilenceUsage: true,
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput := wantsJSON(cmd)
				c := factory()
				id := args[0]

				isOnline, err := c.GetSIPStatus(id)
				if err != nil {
					log.Debugf("Failed to get SIP status: %v", err)
					return failCmd(cmd, err)
				}

				if jsonOutput {
					out := map[string]any{
						"sip":       id,
						"is_online": isOnline,
						"status":    map[bool]string{true: "online", false: "offline"}[isOnline],
					}
					b, _ := json.MarshalIndent(out, "", "  ")
					fmt.Println(string(b))
				} else {
					// Print a small table with SIP ID and status to match CLI style
					printSIPStatusTable(id, isOnline)
				}
				return nil
			},
		},
	)

	return cmd
}

// printSIPTable renders SIP accounts as a formatted table
func printSIPTable(sips []map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	// Print header
	_, _ = fmt.Fprintln(w, "ID\tDISPLAY NAME\tLINES")
	_, _ = fmt.Fprintln(w, "---\t------------\t-----")

	// Print rows
	for _, sip := range sips {
		id := getStringField(sip, "id", "-")
		displayName := getStringField(sip, "display_name", "-")
		lines := getStringField(sip, "lines", "-")

		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\n", id, displayName, lines)
	}
}

// printSIPStatusTable prints the status of a single SIP account in table form
func printSIPStatusTable(id string, isOnline bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	_, _ = fmt.Fprintln(w, "SIP\tSTATUS")
	_, _ = fmt.Fprintln(w, "---\t------")

	status := "offline"
	if isOnline {
		status = "online"
	}

	_, _ = fmt.Fprintf(w, "%s\t%s\n", id, status)
}
