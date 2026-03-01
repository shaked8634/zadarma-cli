package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/zadarma/zadarma-cli/internal/log"
)

// NewPhoneListCmd returns the 'phone list' command
func NewPhoneListCmd(factory ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:   "list [number...]",
		Short: "List phone numbers",
		Long: `List phone numbers (DIDs) owned by your account.

If no numbers are specified, lists all phone numbers.
If one or more numbers are specified, retrieves information for those specific numbers.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonOutput := wantsJSON(cmd)
			c := factory()

			dids, err := c.GetDirectNumbers(args...)
			if err != nil {
				log.Debugf("Failed to get phone numbers: %v", err)
				return failCmd(cmd, err)
			}

			if jsonOutput {
				out, _ := json.MarshalIndent(dids, "", "  ")
				fmt.Println(string(out))
			} else {
				if len(dids) == 0 {
					fmt.Println("No phone numbers found.")
					return nil
				}
				printPhoneTable(dids)
			}
			return nil
		},
	}
}

// printPhoneTable renders owned phone numbers as a formatted table
func printPhoneTable(phones []map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer func() { _ = w.Flush() }()

	// Print header
	fmt.Fprintln(w, "NUMBER\tCOUNTRY\tSTATUS\tDESCRIPTION\tEXPIRES")
	fmt.Fprintln(w, "------\t-------\t------\t-----------\t-------")

	// Print rows
	for _, phone := range phones {
		number := getStringField(phone, "number", "-")
		country := getStringField(phone, "country", "-")
		status := getStringField(phone, "status", "-")
		description := getStringField(phone, "description", "-")
		stopDate := getStringField(phone, "stop_date", "-")

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", number, country, status, description, stopDate)
	}
}
