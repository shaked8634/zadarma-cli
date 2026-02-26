package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewPBXCmd creates the 'pbx' command group.
func NewPBXCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pbx",
		Short: "Manage PBX configuration",
		Long:  "PBX configuration and management commands",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "info",
			Short: "Get PBX information",
			RunE: func(cmd *cobra.Command, args []string) error {
				jsonOutput, _ := cmd.Flags().GetBool("json")
				c := factory()

				pbxInfo, err := c.GetPBXInfo()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}

				if jsonOutput {
					out, _ := json.MarshalIndent(pbxInfo, "", "  ")
					fmt.Println(string(out))
				} else {
					fmt.Printf("PBX Name: %s\n", pbxInfo["name"])
					fmt.Printf("Status: %s\n", pbxInfo["status"])
				}
				return nil
			},
		},
	)

	return cmd
}
