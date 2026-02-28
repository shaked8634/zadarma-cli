package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewDIDCmd creates the 'did' command group.
func NewDIDCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "did",
		Short: "Manage phone numbers (DIDs)",
		Long:  "Phone number (DID) management commands",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List phone numbers",
			RunE: func(cmd *cobra.Command, args []string) error {
				outputFormat, _ := cmd.Root().PersistentFlags().GetString("output")
				c := factory()

				dids, err := c.GetDIDs()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}

				if outputFormat == "json" {
					out, _ := json.MarshalIndent(dids, "", "  ")
					fmt.Println(string(out))
				} else {
					for _, did := range dids {
						fmt.Printf("%s (%s)\n", did["number"], did["type"])
					}
				}
				return nil
			},
		},
	)

	return cmd
}
