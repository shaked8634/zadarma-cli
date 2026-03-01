package commands

import (
	"encoding/json"
	"fmt"

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
				if !jsonOutput {
					if jb, _ := cmd.Root().PersistentFlags().GetBool("json"); jb {
						jsonOutput = true
					}
					of, _ := cmd.Root().PersistentFlags().GetString("output")
					jsonOutput = jsonOutput || of == "json"
				}
				c := factory()

				pbxInfo, err := c.GetPBXInfo()
				if err != nil {
					return failCmd(cmd, err)
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
