package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/log"
)

// NewPhoneCountriesCmd returns the 'phone countries' command
func NewPhoneCountriesCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "countries",
		Short: "List countries with available direct numbers",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			countries, err := c.GetDirectCountries()
			if err != nil {
				log.Debugf("Failed to get countries: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(countries, "", "  ")
				fmt.Println(string(out))
			} else {
				for _, country := range countries {
					fmt.Printf("%v: %v\n", country["id"], country["name"])
				}
			}
			return nil
		},
	}
}
