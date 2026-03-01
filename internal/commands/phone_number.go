package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/log"
)

// NewPhoneNumberCmd returns the 'phone number' command
func NewPhoneNumberCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "number <type> <number>",
		Short: "Get information about a virtual number",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			info, err := c.GetDirectNumber(args[0], args[1])
			if err != nil {
				log.Debugf("Failed to get number info: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(info, "", "  ")
				fmt.Println(string(out))
			} else {
				for key, value := range info {
					fmt.Printf("%s: %v\n", key, value)
				}
			}
			return nil
		},
	}
}
