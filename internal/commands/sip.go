package commands

import (
	"encoding/json"
	"fmt"

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
					return failCmd(cmd, err)
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
		&cobra.Command{
			Use:   "status <ID>",
			Short: "Get SIP account status",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput, _ := cmd.Flags().GetBool("json")
				if !jsonOutput {
					of, _ := cmd.Root().PersistentFlags().GetString("output")
					jsonOutput = of == "json"
				}
				c := factory()
				id := args[0]

				isOnline, err := c.GetSIPStatus(id)
				if err != nil {
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
					status := "offline"
					if isOnline {
						status = "online"
					}
					fmt.Printf("SIP %s is %s\n", id, status)
				}
				return nil
			},
		},
	)

	return cmd
}
