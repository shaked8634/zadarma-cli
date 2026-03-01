package commands

import (
	"github.com/spf13/cobra"
)

// NewPhoneCmd creates the 'phone' command group for managing phone numbers and virtual numbers.
func NewPhoneCmd(factory ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phone",
		Short: "Manage phone numbers (DIDs and virtual numbers)",
		Long:  "Commands for managing owned phone numbers and exploring virtual number availability.",
	}

	cmd.AddCommand(
		NewPhoneListCmd(factory),
		NewPhoneCountriesCmd(factory),
		NewPhoneCountryCmd(factory),
		NewPhoneNumberCmd(factory),
	)

	return cmd
}
