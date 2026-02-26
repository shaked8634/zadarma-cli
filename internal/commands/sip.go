package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewSIPCmd creates the 'sip' command group.
func NewSIPCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sip",
		Short: "Manage SIP accounts",
		Long:  "SIP account management commands",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List SIP accounts",
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput, _ := cmd.Flags().GetBool("json")
				c := factory()

				sips, err := c.GetSIPs()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(sips, "", "  ")
					fmt.Println(string(out))
				} else {
					for _, sip := range sips {
						fmt.Printf("SIP: %s (Status: %s)\n", sip["sip_user"], sip["status"])
					}
				}
				return nil
			},
		},
	)

	return cmd
}
