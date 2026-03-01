package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/log"
)

// NewPhoneCountryCmd returns the 'phone country' command
func NewPhoneCountryCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "country <code>",
		Short: "List available direct numbers in a country",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			countryData, err := c.GetDirectCountry(args[0])
			if err != nil {
				log.Debugf("Failed to get country data: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(countryData, "", "  ")
				fmt.Println(string(out))
			} else {
				if len(countryData) == 0 {
					fmt.Printf("No direct numbers found for %s\n", args[0])
					return nil
				}
				for _, entry := range countryData {
					fmt.Printf("%v: %v\n", entry["id"], entry["name"])
				}
			}
			return nil
		},
	}
}
