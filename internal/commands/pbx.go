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

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Get PBX information",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			pbxID, _ := cmd.Flags().GetString("pbx-id")
			numbers, _ := cmd.Flags().GetString("numbers")

			pbxInfo, err := c.GetPBXInfo(pbxID, numbers)
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
	}
	infoCmd.Flags().String("pbx-id", "", "PBX identifier to query")
	infoCmd.Flags().String("numbers", "", "Comma-separated list of numbers to filter")

	cmd.AddCommand(infoCmd)

	return cmd
}
